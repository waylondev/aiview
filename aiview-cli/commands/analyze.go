package commands

import (
	"fmt"
	"os"
	"path/filepath"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/spf13/cobra"
	"github.com/jackwener/aiview/internal/analyzer"
	"github.com/jackwener/aiview/internal/storage"
)

// NewAnalyzeCmd creates the analyze command.
func NewAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze collected data",
		Long:  "Analyze trends and patterns in collected platform data.",
	}

	cmd.AddCommand(newAnalyzeTrendCmd())
	return cmd
}

func newAnalyzeTrendCmd() *cobra.Command {
	var platform, recordType string
	var days int

	cmd := &cobra.Command{
		Use:   "trend",
		Short: "Analyze data trends over time",
		Long:  "Analyze trends in collected data over a specified time period.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Open storage
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return aiverr.APIError("analyze", fmt.Sprintf("get home dir: %v", err))
			}

			dbPath := filepath.Join(homeDir, ".aiview", "data.db")
			var store storage.Storage

			if _, err := os.Stat(dbPath); err == nil {
				store, err = storage.NewSQLiteStorage(dbPath)
				if err != nil {
					return aiverr.APIError("analyze", fmt.Sprintf("open sqlite: %v", err))
				}
			} else {
				store, err = storage.NewJSONFileStorage(filepath.Join(homeDir, ".aiview", "data"))
				if err != nil {
					return fmt.Errorf("open json storage: %w", err)
				}
			}
			defer store.Close()

			// Analyze trend
			a := analyzer.New(store)
			result, err := a.AnalyzeTrend(platform, recordType, days)
			if err != nil {
				return aiverr.APIError("analyze", fmt.Sprintf("analyze trend: %v", err))
			}

			// Render chart
			chart := analyzer.RenderASCIIChart(result, 60)
			fmt.Print(chart)

			return nil
		},
	}

	cmd.Flags().StringVar(&platform, "platform", "bilibili", "Platform name")
	cmd.Flags().StringVar(&recordType, "type", "hot", "Record type")
	cmd.Flags().IntVar(&days, "days", 7, "Number of days to analyze")

	return cmd
}
