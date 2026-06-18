package douyin

import (
	"fmt"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewStatusCmd creates the status command.
func NewStatusCmd(statusFn func() (map[string]interface{}, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check Douyin login status",
		Long: `Check the current login status for Douyin platform.

Examples:
  aiview douyin status
  aiview douyin status --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format := output.MustGetFormat(cmd)

			result, err := statusFn()
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get status: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			loggedIn := false
			if v, ok := result["logged_in"]; ok {
				loggedIn, _ = v.(bool)
			}

			if loggedIn {
				fmt.Println("✅ Douyin: Logged in")
			} else {
				fmt.Println("❌ Douyin: Not logged in")
			}
			return nil
		},
	}
}