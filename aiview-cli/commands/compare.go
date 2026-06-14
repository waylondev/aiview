package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/spf13/cobra"
	"github.com/jackwener/aiview/internal/analyzer"
	"github.com/jackwener/aiview/internal/storage"
)

// NewCompareCmd creates the compare command.
func NewCompareCmd() *cobra.Command {
	var keyword, platforms string

	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare data across platforms",
		Long:  "Compare the same keyword or topic across multiple platforms.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Open storage
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return aiverr.APIError("compare", fmt.Sprintf("get home dir: %v", err))
			}

			dbPath := filepath.Join(homeDir, ".aiview", "data.db")
			var store storage.Storage

			if _, err := os.Stat(dbPath); err == nil {
				store, err = storage.NewSQLiteStorage(dbPath)
				if err != nil {
					return fmt.Errorf("open sqlite: %w", err)
				}
			} else {
				store, err = storage.NewJSONFileStorage(filepath.Join(homeDir, ".aiview", "data"))
				if err != nil {
					return fmt.Errorf("open json storage: %w", err)
				}
			}
			defer store.Close()

			// Parse platforms
			platformList := strings.Split(platforms, ",")

			// Compare
			a := analyzer.New(store)
			results, err := a.ComparePlatforms(keyword, platformList)
			if err != nil {
				return aiverr.APIError("compare", fmt.Sprintf("compare platforms: %v", err))
			}

			// Render table
			table := analyzer.RenderCompareTable(results, keyword)
			fmt.Print(table)

			return nil
		},
	}

	cmd.Flags().StringVar(&keyword, "keyword", "", "Keyword to compare")
	cmd.Flags().StringVar(&platforms, "platforms", "bilibili,douyin,weibo", "Comma-separated platform list")
	cmd.MarkFlagRequired("keyword")

	return cmd
}
