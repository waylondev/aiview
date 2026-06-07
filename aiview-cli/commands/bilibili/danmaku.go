package bilibili

import (
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
				return output.EmitSuccess(map[string]interface{}{
					"bvid":      bvid,
					"danmaku":   string(xmlData),
					"byte_size": len(xmlData),
				}, format)
			}

			fmt.Printf("💬 Danmaku for %s\n", bvid)
			fmt.Printf("   Raw XML data: %d bytes\n", len(xmlData))
			fmt.Println()
			// Print first few lines if it's text
			content := string(xmlData)
			if len(content) > 500 {
				content = content[:500]
			}
			fmt.Println(content)
			if len(xmlData) > 500 {
				fmt.Println("... (truncated)")
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