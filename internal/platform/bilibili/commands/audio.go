package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewAudioCmd creates the audio command.
func NewAudioCmd(getClient func() Client) *cobra.Command {
	var (
		outputDir string
		noSplit   bool
		segment   int
	)

	cmd := &cobra.Command{
		Use:   "audio <BV>",
		Short: "Download video audio",
		Long:  `Download the audio stream of a Bilibili video.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			audioURL, err := client.GetAudioURL(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get audio stream: %v", err), format)
				return err
			}

			if outputDir == "" {
				outputDir = filepath.Join(os.TempDir(), "aiview", bvid)
			}

			if err := os.MkdirAll(outputDir, 0755); err != nil {
				output.EmitError("internal_error", fmt.Sprintf("Failed to create directory: %v", err), format)
				return err
			}

			outputPath := filepath.Join(outputDir, bvid+".m4a")

			// Download audio
			req, err := http.NewRequest("GET", audioURL, nil)
			if err != nil {
				output.EmitError("network_error", fmt.Sprintf("Failed to create request: %v", err), format)
				return err
			}
			req.Header.Set("User-Agent", "Mozilla/5.0")
			req.Header.Set("Referer", "https://www.bilibili.com")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				output.EmitError("network_error", fmt.Sprintf("Download failed: %v", err), format)
				return err
			}
			defer resp.Body.Close()

			f, err := os.Create(outputPath)
			if err != nil {
				output.EmitError("internal_error", fmt.Sprintf("Failed to create file: %v", err), format)
				return err
			}
			defer f.Close()

			written, err := io.Copy(f, resp.Body)
			if err != nil {
				output.EmitError("network_error", fmt.Sprintf("Download failed: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{
					"path":  outputPath,
					"bytes": written,
				}, format)
			}

			fmt.Printf("✅ Audio downloaded to: %s\n", outputPath)
			fmt.Printf("   Size: %d bytes\n", written)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	cmd.Flags().BoolVar(&noSplit, "no-split", false, "Do not split audio")
	cmd.Flags().IntVar(&segment, "segment", 25, "Seconds per segment")

	return cmd
}