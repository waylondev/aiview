package bilibili

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jackwener/aiview/internal/helper"

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
			format := output.MustGetFormat(cmd)

			uid, err := strconv.Atoi(args[0])
			if err != nil {
				output.EmitError("invalid_input", "UID must be a number", format)
				return err
			}

			result, err := client.GetUserCollections(uid)
			if err != nil {
				if strings.Contains(err.Error(), "-400") {
					output.EmitError("no_collections", "This user has no video collections (合集/系列). Try another UID.", format)
				} else {
					output.EmitError("api_error", fmt.Sprintf("Failed to get collections: %v", err), format)
				}
				return err
			}

			if format == output.FormatJSON || format == output.FormatYAML {
				return output.EmitSuccess(result, format)
			}

			d := helper.GetMap(result, "data")
			itemsLists := helper.GetMap(d, "items_lists")
			seasonsList := helper.GetSlice(itemsLists, "seasons_list")
			seriesList := helper.GetSlice(itemsLists, "series_list")

			fmt.Printf("📂 Collections (UID: %d):\n\n", uid)
			if len(seasonsList) == 0 && len(seriesList) == 0 {
				fmt.Println("  No collections")
				return nil
			}
			count := 0
			for _, item := range seasonsList {
				count++
				m := item.(map[string]interface{})
				name := helper.GetString(m, "title")
				total := helper.GetInt(m, "total")
				fmt.Printf("  %d. [Season] %s (%d videos)\n", count, name, total)
			}
			for _, item := range seriesList {
				count++
				m := item.(map[string]interface{})
				name := helper.GetString(m, "title")
				total := helper.GetInt(m, "total")
				fmt.Printf("  %d. [Series] %s (%d videos)\n", count, name, total)
			}
			return nil
		},
	}
	return cmd
}