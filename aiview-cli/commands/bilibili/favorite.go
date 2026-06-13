package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewFavoriteCmd creates the favorite command.
func NewFavoriteCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var fid int
	var delete bool

	cmd := &cobra.Command{
		Use:   "favorite <BV>",
		Short: "Add or remove video from favorites",
		Long: `Add or remove a video from a favorite folder.

Requires --fid to specify the folder ID.
Use --delete to remove from the folder.

Examples:
  aiview bilibili favorite BV1xx --fid 12345
  aiview bilibili favorite BV1xx --fid 12345 --delete`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.GetFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", err.Error(), format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			action := "favorited"
			if delete {
				if err := client.DelFavorite(bvid, fid); err != nil {
					output.EmitError("api_error", fmt.Sprintf("Failed to remove favorite: %v", err), format)
					return err
				}
				action = "unfavorited"
			} else {
				if err := client.AddFavorite(bvid, fid); err != nil {
					output.EmitError("api_error", fmt.Sprintf("Failed to add favorite: %v", err), format)
					return err
				}
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: fmt.Sprintf("video_%s", action)}, format)
			}

			fmt.Printf("✅ Video %s %s\n", bvid, action)
			return nil
		},
	}

	cmd.Flags().IntVar(&fid, "fid", 0, "Favorite folder ID (required)")
	cmd.MarkFlagRequired("fid")
	cmd.Flags().BoolVar(&delete, "delete", false, "Remove from favorites")
	return cmd
}