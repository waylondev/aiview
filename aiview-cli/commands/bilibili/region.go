package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewRegionCmd creates the region command.
func NewRegionCmd(getClient func() Client) *cobra.Command {
	var (
		page  int
		count int
		sort  string
	)

	cmd := &cobra.Command{
		Use:   "region <rid>",
		Short: "View videos by region",
		Long: `View Bilibili videos by region/category.

Region IDs:
  1 = Animation (MAD·AMV)
  3 = Music
  4 = Game
  5 = Entertainment
  11 = Technology
  13 = Dance
  23 = Movie
  36 = Anime
  119 = Food
  129 = Vlog
  130 = Variety
  155 = Fashion
  160 = Study

Examples:
  aiview bilibili region 1          (Animation)
  aiview bilibili region 3 --sort hot
  aiview bilibili region 11 --pn 2 --ps 10`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			rid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "Region ID must be a number", format)
				return err
			}

			result, err := client.GetRegionVideos(rid, page, count, sort)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get region videos: %v", err), format)
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

			archives := helper.GetSlice(data, "archives")
			if len(archives) == 0 {
				fmt.Printf("No videos found in region %d\n", rid)
				return nil
			}

			fmt.Printf("📂 Region %d Videos:\n\n", rid)
			for i, item := range archives {
				m := item.(map[string]interface{})
				owner := helper.GetMap(m, "owner")
				stat := helper.GetMap(m, "stat")
				bvid := helper.GetString(m, "bvid")
				title := helper.GetString(m, "title")
				author := helper.GetString(owner, "name")
				views := helper.GetInt(stat, "view")

				fmt.Printf("  %d. [%s] %s\n", i+1, bvid, title)
				fmt.Printf("     Uploader: %s  Views: %d\n\n", author, views)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&page, "pn", "p", 1, "Page number")
	cmd.Flags().IntVarP(&count, "ps", "n", 20, "Results per page")
	cmd.Flags().StringVar(&sort, "sort", "", "Sort order (hot/click/pubdate)")
	return cmd
}