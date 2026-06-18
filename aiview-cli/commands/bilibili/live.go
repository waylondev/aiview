package bilibili

import (
	"fmt"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewLiveCmd creates the live command.
func NewLiveCmd(getClient func() Client) *cobra.Command {
	var (
		roomID int
		uid    int
	)

	cmd := &cobra.Command{
		Use:   "live",
		Short: "View live room info",
		Long: `View Bilibili live room information by room ID or user ID.

Examples:
  aiview bilibili live --room 12345
  aiview bilibili live --uid 37737161`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			if roomID == 0 && uid == 0 {
				output.EmitError("invalid_input", "Either --room or --uid is required", format)
				return aiverr.InvalidInput("bilibili", "either --room or --uid is required")
			}

			if roomID == 0 {
				roomID = uid
			}

			result, err := client.GetLiveRoomInfo(roomID)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get live room info for %d: %v", roomID, err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := helper.GetMap(result, "data")
			if data == nil || helper.GetInt(data, "room_id") == 0 {
				fmt.Printf("No live room found for ID %d\n", roomID)
				return nil
			}

			title := helper.GetString(data, "title")
			status := helper.GetInt(data, "live_status")
			online := helper.GetInt(data, "online")
			uname := helper.GetString(data, "uname")
			actualRoomID := helper.GetInt(data, "room_id")

			statusStr := "Offline"
			if status == 1 {
				statusStr = "Live"
			} else if status == 2 {
				statusStr = "Rerunning"
			}

			fmt.Printf("Live Room: %s\n\n", title)
			fmt.Printf("  Room ID:   %d\n", actualRoomID)
			fmt.Printf("  Host:      %s\n", uname)
			fmt.Printf("  Status:    %s\n", statusStr)
			if status == 1 {
				fmt.Printf("  Viewers:   %d\n", online)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&roomID, "room", 0, "Live room ID")
	cmd.Flags().IntVar(&uid, "uid", 0, "User ID (used as room ID lookup)")
	return cmd
}