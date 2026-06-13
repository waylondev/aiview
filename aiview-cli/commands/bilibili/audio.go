package bilibili

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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
		Long:  `Download the audio stream of a Bilibili video. Use --segment to split into WAV segments, or --no-split to keep the full file.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.GetFormat(cmd)

			bvid, err := ExtractBVID(args[0])
			if err != nil {
				output.EmitError("invalid_input", err.Error(), format)
				return err
			}

			// Get video info for title
			info, err := client.GetVideoInfo(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get video info: %v", err), format)
				return err
			}

			audioURL, err := client.GetAudioURL(bvid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get audio stream: %v", err), format)
				return err
			}

			if outputDir == "" {
				outputDir = filepath.Join(os.TempDir(), "aiview", sanitizeTitle(info.Title))
			}

			if err := os.MkdirAll(outputDir, 0755); err != nil {
				output.EmitError("internal_error", fmt.Sprintf("Failed to create directory: %v", err), format)
				return err
			}

			if noSplit {
				// Download full audio without splitting
				outputPath := filepath.Join(outputDir, sanitizeTitle(info.Title)+".m4a")
				nbytes, err := downloadAudio(audioURL, outputPath)
				if err != nil {
					output.EmitError("network_error", fmt.Sprintf("Download failed: %v", err), format)
					return err
				}

				if format == output.FormatJSON || format == output.FormatYAML {
					return output.EmitSuccess(map[string]interface{}{
						"path":  outputPath,
						"bytes": nbytes,
						"size":  float64(nbytes) / (1024 * 1024),
					}, format)
				}

				fmt.Printf("✅ Audio saved to: %s\n", outputPath)
				fmt.Printf("   Size: %.1f MB\n", float64(nbytes)/(1024*1024))
				return nil
			}

			// Download to temp file, then optionally split
			tmpPath := filepath.Join(outputDir, "_raw.m4s")
			nbytes, err := downloadAudio(audioURL, tmpPath)
			if err != nil {
				output.EmitError("network_error", fmt.Sprintf("Download failed: %v", err), format)
				return err
			}

			// Check if ffmpeg is available for splitting
			if !hasFFmpeg() {
				// No ffmpeg, just rename the temp file
				finalPath := filepath.Join(outputDir, sanitizeTitle(info.Title)+".m4a")
				if err := os.Rename(tmpPath, finalPath); err != nil {
					os.Rename(tmpPath, finalPath) // best effort
				}
				if format == output.FormatJSON || format == output.FormatYAML {
					return output.EmitSuccess(map[string]interface{}{
						"path":  finalPath,
						"bytes": nbytes,
						"size":  float64(nbytes) / (1024 * 1024),
					}, format)
				}
				fmt.Printf("✅ Audio saved to: %s\n", finalPath)
				fmt.Printf("   Size: %.1f MB\n", float64(nbytes)/(1024*1024))
				fmt.Println("   (ffmpeg not found, WAV splitting requires ffmpeg)")
				return nil
			}

			// Split with ffmpeg
			segmentStr := strconv.Itoa(segment)
			pattern := filepath.Join(outputDir, "seg_%03d.wav")
			cmdSplit := exec.Command("ffmpeg", "-i", tmpPath,
				"-f", "segment",
				"-segment_time", segmentStr,
				"-ac", "1",     // mono
				"-ar", "16000", // 16kHz
				"-c:a", "pcm_s16le",
				"-y",
				pattern)
			splitOutput, err := cmdSplit.CombinedOutput()
			if err != nil {
				// Clean up temp file
				os.Remove(tmpPath)
				output.EmitError("internal_error", fmt.Sprintf("Failed to split audio with ffmpeg: %v\n%s", err, string(splitOutput)), format)
				return err
			}

			// List segments
			segments, _ := filepath.Glob(filepath.Join(outputDir, "seg_*.wav"))

			// Clean up temp file
			os.Remove(tmpPath)

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(map[string]interface{}{
					"segments":  len(segments),
					"directory": outputDir,
					"duration":  segment,
				}, format)
			}

			fmt.Printf("✅ Split complete: %d segments (each ~%ds)\n", len(segments), segment)
			fmt.Printf("   Output: %s\n", outputDir)
			for _, seg := range segments {
				info, _ := os.Stat(seg)
				sizeKB := float64(info.Size()) / 1024
				fmt.Printf("   %s  (%.0f KB)\n", filepath.Base(seg), sizeKB)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
	cmd.Flags().BoolVar(&noSplit, "no-split", false, "Save full audio without splitting")
	cmd.Flags().IntVar(&segment, "segment", 25, "Seconds per segment (requires ffmpeg)")

	return cmd
}

func downloadAudio(url, outputPath string) (int64, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	f, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	return io.Copy(f, resp.Body)
}

func sanitizeTitle(title string) string {
	r := strings.NewReplacer(
		"<", "_", ">", "_", ":", "_", "\"", "_",
		"/", "_", "\\", "_", "|", "_", "?", "_", "*", "_",
	)
	return strings.TrimSpace(r.Replace(title))
}

func hasFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}