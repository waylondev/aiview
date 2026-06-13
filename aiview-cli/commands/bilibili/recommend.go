package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewRecommendCmd creates the recommend command.
func NewRecommendCmd(getClient func() Client) *cobra.Command {
	var (
		fresh bool
		page  int
		max   int
	)

	cmd := &cobra.Command{
		Use:   "recommend",
		Short: "View recommended videos",
		Long:  `View Bilibili recommended videos on the homepage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			data, err := client.GetRecommendVideos(fresh, page)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get recommendations: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(data, format)
			}

			d := helper.GetMap(data, "data")
			list := helper.GetSlice(d, "item")

			if max > 0 && max < len(list) {
				list = list[:max]
			}

			fmt.Println("🎯 Recommended Videos:")
			for i, item := range list {
				m := item.(map[string]interface{})
				owner := helper.GetMap(m, "owner")
				stat := helper.GetMap(m, "stat")
				title := helper.GetString(m, "title")
				bvid := helper.GetString(m, "bvid")
				name := helper.GetString(owner, "name")
				view := helper.GetInt(stat, "view")

				fmt.Printf("  %d. [%s] %s\n", i+1, bvid, title)
				fmt.Printf("     Uploader: %s  Views: %s\n\n", name, output.FormatCount(view))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&fresh, "fresh", false, "Get fresh recommendations")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	cmd.Flags().IntVarP(&max, "max", "n", 20, "Max results")
	return cmd
}