package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewHotCmd creates the hot command.
func NewHotCmd(getClient func() Client) *cobra.Command {
	var maxResults int

	cmd := &cobra.Command{
		Use:   "hot",
		Short: "View hot videos",
		Long:  `View Bilibili trending videos.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			videos, err := client.GetHotVideos(1, maxResults)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get hot videos: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": videos}, format)
			}

			fmt.Println("🔥 Hot Videos:")
			for i, v := range videos {
				fmt.Printf("  %d. [%s] %s\n", i+1, v.BVID, v.Title)
				fmt.Printf("     Uploader: %s  Views: %s\n\n", v.Owner.Name, output.FormatCount(v.Stats.View))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 20, "Max results")
	return cmd
}

// NewRankCmd creates the rank command.
func NewRankCmd(getClient func() Client) *cobra.Command {
	var maxResults int
	var rid int
	var day int
	var typeStr string

	cmd := &cobra.Command{
		Use:   "rank",
		Short: "View rankings",
		Long:  `View Bilibili rankings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			videos, err := client.GetRankVideos(rid, day, typeStr)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get rankings: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": videos[:min(maxResults, len(videos))]}, format)
			}

			fmt.Println("🏆 Rankings:")
			for i, v := range videos[:min(maxResults, len(videos))] {
				fmt.Printf("  %d. [%s] %s\n", i+1, v.BVID, v.Title)
				fmt.Printf("     Uploader: %s  Views: %s\n\n", v.Owner.Name, output.FormatCount(v.Stats.View))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 20, "Max results")
	cmd.Flags().IntVar(&rid, "rid", 0, "Region ID (0=all)")
	cmd.Flags().IntVar(&day, "day", 3, "Days: 1/3/7/30")
	cmd.Flags().StringVar(&typeStr, "type", "all", "Type: all/origin/rookie")
	return cmd
}

// NewFeedCmd creates the feed command.
func NewFeedCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var maxResults int

	cmd := &cobra.Command{
		Use:   "feed",
		Short: "View feed",
		Long:  `View dynamic feed (login required).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", err.Error(), format)
				return err
			}
			_ = cred

			items, err := client.GetDynamicFeed("")
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get feed: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				limited := items
				if maxResults > 0 && maxResults < len(items) {
					limited = items[:maxResults]
				}
				return output.EmitSuccess(map[string]interface{}{"items": limited}, format)
			}

			fmt.Println("📡 Feed:")
			for i, d := range items {
				if maxResults > 0 && i >= maxResults {
					break
				}
				fmt.Printf("  %s\n", d.Author)
				if d.Text != "" {
					text := d.Text
					if len(text) > 100 {
						text = text[:100]
					}
					fmt.Printf("  %s\n\n", text)
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 20, "Max results")
	return cmd
}