package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewOnlineCmd creates the online command.
func NewOnlineCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "online <BV>",
		Short: "View real-time video online viewers",
		Long: `View the real-time number of people watching a Bilibili video.

Shows total online viewers and web platform viewers.

Examples:
  aiview bilibili online BV1xx411c7m9`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			result, err := client.GetVideoOnlineCount(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get online count: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := helper.GetMap(result, "data")
			if data == nil {
				fmt.Println("No data returned")
				return nil
			}

			total := helper.GetInt(data, "total")
			webOnline := helper.GetInt(data, "count")

			fmt.Printf("📺 Online Viewers: %s\n\n", bvid)
			fmt.Printf("  Total viewers:   %d\n", total)
			fmt.Printf("  Web viewers:     %d\n", webOnline)
			return nil
		},
	}
}