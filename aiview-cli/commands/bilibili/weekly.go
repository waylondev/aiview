package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewWeeklyCmd creates the weekly command.
func NewWeeklyCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "weekly <number>",
		Short: "View weekly hot video series",
		Long: `View Bilibili's "每周必看" (weekly must-watch) video series by week number.

Examples:
  aiview bilibili weekly 244
  aiview bilibili weekly 300`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			number, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "Week number must be a number", format)
				return err
			}

			result, err := client.GetWeeklyHotVideos(number)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get weekly videos: %v", err), format)
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

			list := helper.GetSlice(data, "list")
			if len(list) == 0 {
				fmt.Printf("No videos found for week %d\n", number)
				return nil
			}

			fmt.Printf("📅 Weekly Hot Videos (Week %d):\n\n", number)
			for i, item := range list {
				m := item.(map[string]interface{})
				owner := helper.GetMap(m, "owner")
				stat := helper.GetMap(m, "stat")
				bvid := helper.GetString(m, "bvid")
				vtitle := helper.GetString(m, "title")
				author := helper.GetString(owner, "name")
				views := helper.GetInt(stat, "view")

				fmt.Printf("  %d. [%s] %s\n", i+1, bvid, vtitle)
				fmt.Printf("     UP主: %s  Views: %d\n\n", author, views)
			}
			return nil
		},
	}
}