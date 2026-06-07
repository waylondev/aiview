package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewCommentCmd creates the comment command.
func NewCommentCmd(getClient func() Client) *cobra.Command {
	var cursor int

	cmd := &cobra.Command{
		Use:   "comment <video_id>",
		Short: "View Douyin video comments",
		Long: `View comments of a Douyin video by video ID.

Requires login cookie for full access.

Examples:
  aiview douyin comment 123456789
  aiview douyin comment 123456789 --cursor 20`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			result, err := client.GetVideoComments(args[0], cursor)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get comments: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := getMap(result, "data")
			if data == nil {
				fmt.Println("No comments data returned")
				return nil
			}

			comments := getSlice(data, "comments")
			if len(comments) == 0 {
				fmt.Println("No comments found")
				return nil
			}

			fmt.Printf("💬 Douyin Comments (video: %s):\n\n", args[0])
			for i, item := range comments {
				m := item.(map[string]interface{})
				user := getMap(m, "user")
				username := getString(user, "nickname")
				text := getString(m, "text")

				fmt.Printf("  %2d. ", i+1)
				if username != "" {
					fmt.Printf("%s: ", username)
				}
				fmt.Printf("%s\n", text)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&cursor, "cursor", "c", 0, "Comment cursor/pagination offset")
	return cmd
}