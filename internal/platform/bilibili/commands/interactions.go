package commands

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLikeCmd creates the like command.
func NewLikeCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "like <BV号>",
		Short: "点赞视频",
		Long:  `点赞视频（需要登录且有写权限）。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录并具有写权限，请使用 aiview bilibili login 登录", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			if err := client.LikeVideo(bvid, false); err != nil {
				output.EmitError("api_error", fmt.Sprintf("点赞失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "like"}, format)
			}

			fmt.Println("✅ 已点赞")
			return nil
		},
	}
}

// NewCoinCmd creates the coin command.
func NewCoinCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var num int

	cmd := &cobra.Command{
		Use:   "coin <BV号>",
		Short: "投币视频",
		Long:  `投币视频（需要登录且有写权限）。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录并具有写权限", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			if err := client.CoinVideo(bvid, num); err != nil {
				output.EmitError("api_error", fmt.Sprintf("投币失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "coin"}, format)
			}

			fmt.Printf("✅ 已投 %d 个硬币\n", num)
			return nil
		},
	}

	cmd.Flags().IntVarP(&num, "num", "n", 1, "投币数量 (1-2)")
	return cmd
}

// NewTripleCmd creates the triple command.
func NewTripleCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "triple <BV号>",
		Short: "一键三连",
		Long:  `一键三连（点赞+投币+收藏）（需要登录且有写权限）。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录并具有写权限", format)
				return err
			}
			_ = cred

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			if err := client.TripleVideo(bvid); err != nil {
				output.EmitError("api_error", fmt.Sprintf("三连失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "triple"}, format)
			}

			fmt.Println("✅ 一键三连成功")
			return nil
		},
	}
}

// NewUnfollowCmd creates the unfollow command.
func NewUnfollowCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "unfollow <UID>",
		Short: "取消关注",
		Long:  `取消关注用户（需要登录且有写权限）。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			cred, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", "需要登录并具有写权限", format)
				return err
			}
			_ = cred

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID 必须是数字", format)
				return err
			}

			if err := client.UnfollowUser(uid); err != nil {
				output.EmitError("api_error", fmt.Sprintf("取消关注失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "unfollow"}, format)
			}

			fmt.Println("✅ 已取消关注")
			return nil
		},
	}
}