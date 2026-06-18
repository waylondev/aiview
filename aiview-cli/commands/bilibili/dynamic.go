package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/helper"

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
			format := output.MustGetFormat(cmd)

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

			d := helper.GetMap(result, "data")
			items := helper.GetSlice(d, "items")

			fmt.Printf("📡 Dynamics (UID: %d):\n\n", uid)
			if len(items) == 0 {
				fmt.Println("  No dynamics")
				return nil
			}
			for i, item := range items {
				m := item.(map[string]interface{})
				modules := helper.GetMap(m, "modules")
				author := helper.GetMap(modules, "module_author")
				desc := helper.GetMap(helper.GetMap(modules, "module_dynamic"), "desc")

				name := helper.GetString(author, "name")
				text := helper.GetString(desc, "text")
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

// NewDynamicPostCmd creates the dynamic-post command.
func NewDynamicPostCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "dynamic-post <text>",
		Short: "Post a text dynamic",
		Long: `Post a plain text dynamic to your Bilibili feed (login and write permission required).

Examples:
  aiview bilibili dynamic-post "Hello World!"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			_, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", err.Error(), format)
				return err
			}

			result, err := client.PostDynamic(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to post dynamic: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			fmt.Println("✅ Dynamic posted successfully")
			return nil
		},
	}
}

// NewDynamicDeleteCmd creates the dynamic-delete command.
func NewDynamicDeleteCmd(authStore AuthProvider, getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "dynamic-delete <id>",
		Short: "Delete a dynamic",
		Long: `Delete a dynamic by its ID (login and write permission required).

Examples:
  aiview bilibili dynamic-delete 123456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

			_, err := authStore.RequireCredential(true)
			if err != nil {
				output.EmitError("not_authenticated", err.Error(), format)
				return err
			}

			dynamicID, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "Dynamic ID must be a number", format)
				return err
			}

			result, err := client.DeleteDynamic(dynamicID)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to delete dynamic: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			fmt.Println("✅ Dynamic deleted successfully")
			return nil
		},
	}
}