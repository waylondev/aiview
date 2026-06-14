package commands

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/storage"
	"github.com/spf13/cobra"
)

// NewExportCmd creates the export command.
func NewExportCmd() *cobra.Command {
	var platform, recordType, format, output string
	var limit int

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export collected data",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Open storage
			homeDir, _ := os.UserHomeDir()
			var store storage.Storage
			var err error

			dbPath := homeDir + "/.aiview/data.db"
			if _, err := os.Stat(dbPath); err == nil {
				store, err = storage.NewSQLiteStorage(dbPath)
				if err != nil {
					return err
				}
			} else {
				store, err = storage.NewJSONFileStorage(homeDir + "/.aiview/data")
				if err != nil {
					return err
				}
			}
			defer store.Close()

			records, err := store.Query(platform, recordType, limit)
			if err != nil {
				return err
			}

			switch format {
			case "json":
				data, _ := json.MarshalIndent(records, "", "  ")
				if output != "" {
					return os.WriteFile(output, data, 0644)
				}
				fmt.Println(string(data))
			case "csv":
				var w *csv.Writer
				if output != "" {
					f, err := os.Create(output)
					if err != nil {
						return err
					}
					defer f.Close()
					w = csv.NewWriter(f)
				} else {
					w = csv.NewWriter(os.Stdout)
				}
				for _, r := range records {
					data, _ := json.Marshal(r.Data)
					w.Write([]string{r.Platform, r.Type, string(data), r.CollectedAt.String()})
				}
				w.Flush()
			default:
				return aiverr.InvalidInput("export", fmt.Sprintf("unknown format: %s", format))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&platform, "platform", "bilibili", "Platform name")
	cmd.Flags().StringVar(&recordType, "type", "hot", "Record type")
	cmd.Flags().StringVar(&format, "format", "json", "Output format: json or csv")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	cmd.Flags().IntVar(&limit, "limit", 100, "Max records to export")
	return cmd
}
