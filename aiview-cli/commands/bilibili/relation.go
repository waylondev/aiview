package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewRelationCmd creates the relation command.
func NewRelationCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "relation <UID>",
		Short: "View user relation status",
		Long: `View the relation status between the current user and another user.

Requires login (read-only permission is sufficient).

Examples:
  aiview bilibili relation 12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			result, err := client.GetRelationStat(uid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get relation: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			data := helper.GetMap(result, "data")
			if data == nil {
				fmt.Println("No data returned")
				return nil
			}

			following := helper.GetInt(data, "following") == 1
			follower := helper.GetInt(data, "follower") == 1

			fmt.Printf("👥 Relation with UID %d:\n\n", uid)
			fmt.Printf("  You follow them:  %s\n", boolMark(following))
			fmt.Printf("  They follow you:  %s\n", boolMark(follower))
			return nil
		},
	}
}

func boolMark(b bool) string {
	if b {
		return "✅ Yes"
	}
	return "❌ No"
}