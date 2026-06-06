package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewCommentCmd creates the comment command.
func NewCommentCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var (
		message string
		root    int
		parent  int
		page    int
		sort    int
	)

	cmd := &cobra.Command{
		Use:   "comment <BV or URL>",
		Short: "View or post comments",
		Long:  `View comments on a video, or post/delete comments (login and write permission required for post/delete).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			info, err := client.GetVideoInfo(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			// Post mode
			if message != "" {
				cred, err := authStore.RequireCredential(true)
				if err != nil {
					output.EmitError("not_authenticated", "Login with write permission required to post comments", format)
					return err
				}
				_ = cred

				if err := client.PostComment(info.AID, message, root, parent); err != nil {
					output.EmitError("api_error", fmt.Sprintf("Failed to post comment: %v", err), format)
					return err
				}

				if format == output.FormatJSON || format == output.FormatYAML {
					return output.EmitSuccess(ActionResult{Success: true, Action: "comment_post"}, format)
				}
				fmt.Println("✅ Comment posted")
				return nil
			}

			// View comments mode
			data, err := client.GetVideoCommentsRaw(info.AID, page, sort)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get comments: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(data, format)
			}

			// Parse and display
			d := getMap(data, "data")
			replies := getSlice(d, "replies")

			if sort == 2 {
				fmt.Println("💬 Hot Comments:")
			} else {
				fmt.Println("💬 Comments:")
			}

			if len(replies) == 0 {
				fmt.Println("  No comments")
				return nil
			}

			for i, r := range replies {
				m := r.(map[string]interface{})
				member := getMap(m, "member")
				content := getMap(m, "content")
				name := getString(member, "uname")
				msg := getString(content, "message")
				like := getInt(m, "like")
				rpid := getInt(m, "rpid")

				fmt.Printf("  %d. %s (👍 %d)\n", i+1, name, like)
				fmt.Printf("     ID: %d\n", rpid)
				if len(msg) > 150 {
					msg = msg[:150]
				}
				fmt.Printf("     %s\n\n", msg)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Post a comment with the given message")
	cmd.Flags().IntVar(&root, "root", 0, "Root comment ID for reply")
	cmd.Flags().IntVar(&parent, "parent", 0, "Parent comment ID for reply")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	cmd.Flags().IntVar(&sort, "sort", 0, "Sort order: 0=time, 2=hot")
	return cmd
}

// NewCommentDeleteCmd creates the comment delete command.
func NewCommentDeleteCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "comment-delete <BV> <RPID>",
		Short: "Delete a comment",
		Long:  `Delete your own comment on a video (login and write permission required).`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "Login with write permission required", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			rpid, err := strconv.Atoi(args[1])
			if err != nil {
				output.EmitError("invalid_input", "RPID must be a number", format)
				return err
			}

			info, err := client.GetVideoInfo(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			if err := client.DeleteComment(info.AID, rpid); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to delete comment: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "comment_delete"}, format)
			}

			fmt.Println("✅ Comment deleted")
			return nil
		},
	}
}