package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewDynamicCmd creates the dynamic command.
func NewDynamicCmd(getClient func() Client) *cobra.Command {
	var page int

	cmd := &cobra.Command{
		Use:   "dynamic <UID>",
		Short: "View user's dynamics",
		Long: `View a user's dynamic feed from their personal space.

Examples:
  aiview bilibili dynamic 12345
  aiview bilibili dynamic 12345 --page 2`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			result, err := client.GetUserDynamics(uid, page)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get dynamics: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			d := getMap(result, "data")
			items := getSlice(d, "items")

			fmt.Printf("📡 Dynamics (UID: %d):\n\n", uid)
			if len(items) == 0 {
				fmt.Println("  No dynamics")
				return nil
			}
			for i, item := range items {
				m := item.(map[string]interface{})
				modules := getMap(m, "modules")
				author := getMap(modules, "module_author")
				desc := getMap(getMap(modules, "module_dynamic"), "desc")

				name := getString(author, "name")
				text := getString(desc, "text")
				if len(text) > 120 {
					text = text[:120] + "..."
				}

				fmt.Printf("  %d. %s\n", i+1, name)
				if text != "" {
					fmt.Printf("     %s\n", text)
				}
				fmt.Println()
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number")
	return cmd
}