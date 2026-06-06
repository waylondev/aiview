package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewFavoritesCmd creates the favorites command.
func NewFavoritesCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var max int
	var page int

	cmd := &cobra.Command{
		Use:   "favorites",
		Short: "View favorites",
		Long:  `View favorite folders (login required).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "Login required, use aiview bilibili login", format)
				return err
			}
			_ = cred

			info, err := client.GetSelfInfo()
			if err != nil {
				output.EmitError("not_authenticated", fmt.Sprintf("Failed to get user info: %v", err), format)
				return err
			}

			folders, err := client.GetFavoriteList(info.MID, page)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get favorites: %v", err), format)
				return err
			}

			if max > 0 && max < len(folders) {
				folders = folders[:max]
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": folders}, format)
			}

			fmt.Println("📁 Favorite Folders:")
			for _, f := range folders {
				fmt.Printf("  [%d] %s (%d videos)\n", f.ID, f.Title, f.MediaCount)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&max, "max", "n", 0, "Maximum number of folders to show")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	return cmd
}

// NewFollowingCmd creates the following command.
func NewFollowingCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var max int
	var page int

	cmd := &cobra.Command{
		Use:   "following",
		Short: "View following",
		Long:  `View following list (login required).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "Login required, use aiview bilibili login", format)
				return err
			}
			_ = cred

			info, err := client.GetSelfInfo()
			if err != nil {
				output.EmitError("not_authenticated", fmt.Sprintf("Failed to get user info: %v", err), format)
				return err
			}

			users, err := client.GetFollowingList(info.MID, page)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get following list: %v", err), format)
				return err
			}

			if max > 0 && max < len(users) {
				users = users[:max]
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": users}, format)
			}

			fmt.Println("👥 Following:")
			for _, u := range users {
				fmt.Printf("  %s (UID: %d)\n", u.Name, u.MID)
				if u.Sign != "" {
					fmt.Printf("    %s\n", u.Sign)
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&max, "max", "n", 0, "Maximum number of users to show")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	return cmd
}

// NewHistoryCmd creates the history command.
func NewHistoryCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var max int
	var page int

	cmd := &cobra.Command{
		Use:   "history",
		Short: "View history",
		Long:  `View watch history (login required).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "Login required, use aiview bilibili login", format)
				return err
			}
			_ = cred

			count := 20
			if max > 0 {
				count = max
			}
			items, err := client.GetWatchHistory(page, count)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get watch history: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": items}, format)
			}

			fmt.Println("📺 Watch History:")
			for i, h := range items {
				fmt.Printf("  %d. [%s] %s\n", i+1, h.BVID, h.Title)
				fmt.Printf("     Uploader: %s\n", h.Author)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&max, "max", "n", 0, "Maximum number of history items to show")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	return cmd
}

// NewWatchLaterCmd creates the watch-later command.
func NewWatchLaterCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var max int

	cmd := &cobra.Command{
		Use:   "watch-later",
		Short: "View watch later",
		Long:  `View watch later list (login required).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "Login required, use aiview bilibili login", format)
				return err
			}
			_ = cred

			items, err := client.GetWatchLater()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get watch later list: %v", err), format)
				return err
			}

			if max > 0 && max < len(items) {
				items = items[:max]
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": items}, format)
			}

			fmt.Println("⏰ Watch Later:")
			for i, w := range items {
				fmt.Printf("  %d. [%s] %s\n", i+1, w.BVID, w.Title)
				fmt.Printf("     Uploader: %s  Duration: %s\n", w.Author, w.Duration)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&max, "max", "n", 0, "Maximum number of items to show")
	return cmd
}