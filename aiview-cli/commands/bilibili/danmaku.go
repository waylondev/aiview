package bilibili

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewDanmakuCmd creates the danmaku command.
func NewDanmakuCmd(getClient func() Client) *cobra.Command {
	var (
		outputDir string
		progress  int
	)

	cmd := &cobra.Command{
		Use:   "danmaku <BV>",
		Short: "View or send danmaku",
		Long:  `View danmaku (bullet comments) of a video, or send a danmaku (login and write permission required for send).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.GetFormat(cmd)

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

			// Get the first page's cid
			// We need the cid for danmaku - get it from the video info or a separate call
			// For simplicity, download the danmaku XML for the aid (oid)
			xmlData, err := client.GetVideoDanmaku(info.CID)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get danmaku: %v", err), format)
				return err
			}

			if outputDir != "" {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					output.EmitError("internal_error", fmt.Sprintf("Failed to create directory: %v", err), format)
					return err
				}
				outputPath := filepath.Join(outputDir, bvid+"_danmaku.xml")
				if err := os.WriteFile(outputPath, xmlData, 0644); err != nil {
					output.EmitError("internal_error", fmt.Sprintf("Failed to write danmaku file: %v", err), format)
					return err
				}

				if format == output.FormatJSON || format == output.FormatYAML {
					return output.EmitSuccess(map[string]interface{}{
						"path":  outputPath,
						"bytes": len(xmlData),
					}, format)
				}

				fmt.Printf("✅ Danmaku saved to: %s\n", outputPath)
				fmt.Printf("   Size: %d bytes\n", len(xmlData))
				return nil
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				// protobuf 格式需要专门的解码器，当前以 hex dump 形式展示
				hexDump := hex.Dump(xmlData)
				return output.EmitSuccess(map[string]interface{}{
					"bvid":      bvid,
					"danmaku":   hexDump,
					"byte_size": len(xmlData),
				}, format)
			}

			fmt.Printf("💬 Danmaku for %s\n", bvid)
			fmt.Printf("   Raw protobuf data: %d bytes\n", len(xmlData))
			fmt.Println()
			// protobuf 格式需要专门的解码器，当前以 hex dump 形式展示
			hexDump := hex.Dump(xmlData)
			if len(hexDump) > 1000 {
				fmt.Println(hexDump[:1000])
				fmt.Println("... (truncated)")
			} else {
				fmt.Println(hexDump)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Save danmaku to directory")
	cmd.Flags().IntVar(&progress, "progress", 0, "Video progress in seconds for danmaku position")
	return cmd
}

// NewDanmakuSendCmd creates the danmaku send command.
func NewDanmakuSendCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	var progress int

	cmd := &cobra.Command{
		Use:   "danmaku-send <BV> <message>",
		Short: "Send a danmaku",
		Long:  `Send a danmaku (bullet comment) on a video (login and write permission required).`,
		Args:  cobra.ExactArgs(2),
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

			message := args[1]

			info, err := client.GetVideoInfo(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			if err := client.PostDanmaku(info.AID, info.CID, message, progress); err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to send danmaku: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(ActionResult{Success: true, Action: "danmaku_send"}, format)
			}

			fmt.Println("✅ Danmaku sent")
			return nil
		},
	}

	cmd.Flags().IntVar(&progress, "progress", 0, "Video progress in seconds")
	return cmd
}