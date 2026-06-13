package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewFansCmd creates the fans command.
func NewFansCmd(getClient func() Client) *cobra.Command {
	var page int

	cmd := &cobra.Command{
		Use:   "fans <UID>",
		Short: "View fans list",
		Long:  `View the fans/followers list of a Bilibili user.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			data, err := client.GetFansList(uid, page)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get fans list: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(data, format)
			}

			d := helper.GetMap(data, "data")
			list := helper.GetSlice(d, "list")

			fmt.Printf("👥 Fans (UID: %d):\n\n", uid)
			if len(list) == 0 {
				fmt.Println("  No fans")
				return nil
			}
			for i, u := range list {
				m := u.(map[string]interface{})
				name := helper.GetString(m, "uname")
				mid := helper.GetInt(m, "mid")
				sign := helper.GetString(m, "sign")
				fmt.Printf("  %d. %s (UID: %d)\n", i+1, name, mid)
				if sign != "" {
					fmt.Printf("     %s\n", sign)
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	return cmd
}