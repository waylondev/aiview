package bilibili

import (
	"fmt"

	"github.com/jackwener/aiview/internal/helper"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewSuggestCmd creates the suggest command.
func NewSuggestCmd(getClient func() Client) *cobra.Command {
	return &cobra.Command{
		Use:   "suggest <keyword>",
		Short: "Get search suggestions",
		Long:  `Get Bilibili search suggestions for a keyword.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			data, err := client.SearchSuggest(args[0])
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get suggestions: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(data, format)
			}

			d := helper.GetMap(data, "data")
			showName := helper.GetString(d, "show_name")

			fmt.Printf("🔍 Suggestions for \"%s\":\n\n", args[0])
			if showName == "" {
				fmt.Println("  No suggestions")
				return nil
			}
			fmt.Printf("  %s\n", showName)
			return nil
		},
	}
}