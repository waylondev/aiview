package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserPostsCmd creates the user-posts command.
func NewUserPostsCmd(client Client) *cobra.Command {
	var cursor int

	cmd := &cobra.Command{
		Use:   "user-posts <uid>",
		Short: "View Douyin user's video posts",
		Long: `View video posts of a Douyin user by user ID.

Requires login cookie for full access.

Examples:
  aiview douyin user-posts 123456789
  aiview douyin user-posts 123456789 --cursor 20`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.GetFormat(cmd)

			result, err := client.GetUserPosts(args[0], cursor)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get user posts: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := helper.GetMap(result, "data")
			if data == nil {
				fmt.Println("No posts data returned")
				return nil
			}

			posts := helper.GetSlice(data, "posts")
			if len(posts) == 0 {
				fmt.Println("No posts found")
				return nil
			}

			fmt.Printf("📹 Douyin User Posts (uid: %s):\n\n", args[0])
			for i, item := range posts {
				m := item.(map[string]interface{})
				desc := helper.GetString(m, "desc")
				fmt.Printf("  %2d. %s\n", i+1, desc)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&cursor, "cursor", "c", 0, "Post cursor/pagination offset")
	return cmd
}