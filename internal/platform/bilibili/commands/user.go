package commands

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewUserCmd creates the user command.
func NewUserCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "user <UID>",
		Short: "查看用户信息",
		Long:  `查看 Bilibili 用户信息。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID 必须是数字", format)
				return err
			}

			info, err := client.GetUserInfo(uid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取用户信息失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(info, format)
			}

			fmt.Printf("👤 %s (UID: %d)\n", info.Name, info.MID)
			fmt.Printf("  等级: %d\n", info.Level)
			fmt.Printf("  硬币: %d\n", info.Coins)
			if info.Sign != "" {
				fmt.Printf("  签名: %s\n", info.Sign)
			}
			return nil
		},
	}
}

// NewUserVideosCmd creates the user-videos command.
func NewUserVideosCmd(getClient func() Client) *cobra.Command {
	var maxResults int

	cmd := &cobra.Command{
		Use:   "user-videos <UID>",
		Short: "查看用户视频列表",
		Long:  `查看 Bilibili 用户的视频列表。`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID 必须是数字", format)
				return err
			}

			videos, err := client.GetUserVideos(uid, maxResults)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("获取用户视频失败: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				type result struct {
					Items []VideoInfo `json:"items"`
				}
				return output.EmitSuccess(result{Items: videos}, format)
			}

			fmt.Printf("📹 用户视频 (UID: %d):\n\n", uid)
			for i, v := range videos {
				fmt.Printf("  %d. [%s] %s\n", i+1, v.BVID, v.Title)
				fmt.Printf("     播放: %s  时长: %s\n\n", output.FormatCount(v.Stats.View), v.DurationStr)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&maxResults, "max", "n", 10, "最大结果数")
	return cmd
}