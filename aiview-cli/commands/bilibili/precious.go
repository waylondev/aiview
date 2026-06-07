package bilibili

import (
	"fmt"

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

			data := getMap(result, "data")
			if data == nil {
				fmt.Println("No data returned")
				return nil
			}

			title := getString(data, "title")
			explain := getString(data, "explain")
			list := getSlice(data, "list")

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
				owner := getMap(m, "owner")
				stat := getMap(m, "stat")
				bvid := getString(m, "bvid")
				vtitle := getString(m, "title")
				author := getString(owner, "name")
				views := getInt(stat, "view")
				achievement := getString(m, "achievement")

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