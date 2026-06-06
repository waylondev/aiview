package bilibili

import (
	"fmt"
	"strconv"

	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// NewCollectionCmd creates the collection command.
func NewCollectionCmd(getClient func() Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection <UID>",
		Short: "View user's video collections",
		Long: `View a user's video collection (series) list.

A collection (合集) is a manually curated series of videos by a creator.

Examples:
  aiview bilibili collection 12345`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := GetOutputFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			result, err := client.GetUserCollections(uid)
			if err != nil {
				output.EmitError("api_error", fmt.Sprintf("Failed to get collections: %v", err), format)
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			d := getMap(result, "data")
			list := getSlice(d, "list")

			fmt.Printf("📂 Collections (UID: %d):\n\n", uid)
			if len(list) == 0 {
				fmt.Println("  No collections")
				return nil
			}
			for i, item := range list {
				m := item.(map[string]interface{})
				name := getString(m, "name")
				total := getInt(m, "total")
				fmt.Printf("  %d. %s (%d videos)\n", i+1, name, total)
			}
			return nil
		},
	}
	return cmd
}