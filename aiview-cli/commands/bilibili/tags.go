package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewTagsCmd creates the tags command.
func NewTagsCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "tags <BV or URL>",
		Short: "View video tags",
		Long:  `View tags of a Bilibili video.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			data, err := client.GetVideoTags(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get tags: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(data, format)
			}

			d := helper.GetMap(data, "data")
			list := helper.GetSlice(d, "tags")

			fmt.Printf("🏷️  Tags:\n\n")
			if len(list) == 0 {
				fmt.Println("  No tags")
				return nil
			}
			for _, t := range list {
				m := t.(map[string]interface{})
				name := helper.GetString(m, "tag_name")
				likes := helper.GetInt(m, "likes")
				fmt.Printf("  #%s (👍 %d)\n", name, likes)
			}
			return nil
		},
	}
}