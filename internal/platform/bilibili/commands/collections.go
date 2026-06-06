package commands

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewFavoritesCmd creates the favorites command.
func NewFavoritesCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "favorites",
		Short: "查看收藏夹",
		Long:  `查看收藏夹列表（需要登录）。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录，请使用 aiview bilibili login 登录", format)
				return err
			}
			_ = cred

			info, err := client.GetSelfInfo()
			if err != nil {
				output.EmitError("not_authenticated", fmt.Sprintf("获取用户信息失败: %v", err), format)
				return err
			}

			folders, err := client.GetFavoriteList(info.MID)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取收藏夹失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": folders}, format)
			}

			fmt.Println("📁 收藏夹:")
			for _, f := range folders {
				fmt.Printf("  [%d] %s (%d 个视频)\n", f.ID, f.Title, f.MediaCount)
			}
			return nil
		},
	}
}

// NewFollowingCmd creates the following command.
func NewFollowingCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "following",
		Short: "查看关注列表",
		Long:  `查看关注列表（需要登录）。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录，请使用 aiview bilibili login 登录", format)
				return err
			}
			_ = cred

			info, err := client.GetSelfInfo()
			if err != nil {
				output.EmitError("not_authenticated", fmt.Sprintf("获取用户信息失败: %v", err), format)
				return err
			}

			users, err := client.GetFollowingList(info.MID, 1)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取关注列表失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": users}, format)
			}

			fmt.Println("👥 关注列表:")
			for _, u := range users {
				fmt.Printf("  %s (UID: %d)\n", u.Name, u.MID)
				if u.Sign != "" {
					fmt.Printf("    %s\n", u.Sign)
				}
			}
			return nil
		},
	}
}

// NewHistoryCmd creates the history command.
func NewHistoryCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "查看观看历史",
		Long:  `查看观看历史（需要登录）。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录，请使用 aiview bilibili login 登录", format)
				return err
			}
			_ = cred

			items, err := client.GetWatchHistory(1, 20)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取观看历史失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": items}, format)
			}

			fmt.Println("📺 观看历史:")
			for i, h := range items {
				fmt.Printf("  %d. [%s] %s\n", i+1, h.BVID, h.Title)
				fmt.Printf("     UP主: %s\n", h.Author)
			}
			return nil
		},
	}
}

// NewWatchLaterCmd creates the watch-later command.
func NewWatchLaterCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "watch-later",
		Short: "查看稍后再看",
		Long:  `查看稍后再看列表（需要登录）。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(false)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录，请使用 aiview bilibili login 登录", format)
				return err
			}
			_ = cred

			items, err := client.GetWatchLater()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取稍后再看失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{"items": items}, format)
			}

			fmt.Println("⏰ 稍后再看:")
			for i, w := range items {
				fmt.Printf("  %d. [%s] %s\n", i+1, w.BVID, w.Title)
				fmt.Printf("     UP主: %s  时长: %s\n", w.Author, w.Duration)
			}
			return nil
		},
	}
}