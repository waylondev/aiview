package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewPreciousCmd creates the precious command.
func NewPreciousCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "precious",
		Short: "View must-watch (precious) videos",
		Long: `View Bilibili's curated "入站必刷" (must-watch) classic video collection.

These are hand-picked videos that Bilibili considers essential viewing.

Examples:
  aiview bilibili precious`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			result, err := client.GetPreciousVideos()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get precious videos: %v", err), format)
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

			title := helper.GetString(data, "title")
			explain := helper.GetString(data, "explain")
			list := helper.GetSlice(data, "list")

			fmt.Printf("🏆 %s\n", title)
			if explain != "" {
				fmt.Printf("   %s\n\n", explain)
			}
			if len(list) == 0 {
				fmt.Println("No precious videos found")
				return nil
			}

			for i, item := range list {
				m := item.(map[string]interface{})
				owner := helper.GetMap(m, "owner")
				stat := helper.GetMap(m, "stat")
				bvid := helper.GetString(m, "bvid")
				vtitle := helper.GetString(m, "title")
				author := helper.GetString(owner, "name")
				views := helper.GetInt(stat, "view")
				achievement := helper.GetString(m, "achievement")

				fmt.Printf("  %d. [%s] %s\n", i+1, bvid, vtitle)
				fmt.Printf("     UP主: %s  Views: %d\n", author, views)
				if achievement != "" {
					fmt.Printf("     🏅 %s\n", achievement)
				}
				fmt.Println()
			}
			return nil
		},
	}
}