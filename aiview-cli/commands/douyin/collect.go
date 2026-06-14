package douyin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/auth"
	"github.com/jackwener/aiview/internal/output"
	"github.com/jackwener/aiview/internal/pipeline"
	"github.com/jackwener/aiview/internal/storage"
	"github.com/spf13/cobra"
)

// douyinCollector wraps the Client to implement pipeline.Collector
type douyinCollector struct {
	client Client
}

func (dc *douyinCollector) Collect(recordType string) (map[string]interface{}, error) {
	switch recordType {
	case "hot":
		return dc.client.GetHotSearch()
	case "trending":
		return dc.client.GetTrending()
	case "search":
		return dc.client.Search("热门", 1, 20)
	default:
		return nil, fmt.Errorf("unknown type: %s", recordType)
	}
}

// NewCollectCmd creates the collect command for douyin.
func NewCollectCmd(getClient func() Client, isLoggedIn func() bool) *cobra.Command {
	var types string
	var storeType string

	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Batch collect and store data",
		Long: `Batch collect and store data from Douyin.

Requires login cookie for full access.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.RequireAuth("douyin", isLoggedIn); err != nil {
				return err
			}

			client := getClient()
			format := output.GetFormat(cmd)

			// Create storage
			homeDir, _ := os.UserHomeDir()
			storeDir := filepath.Join(homeDir, ".aiview", "data")

			var store storage.Storage
			var err error
			switch storeType {
			case "sqlite":
				dbPath := filepath.Join(homeDir, ".aiview", "data.db")
				store, err = storage.NewSQLiteStorage(dbPath)
			case "json":
				store, err = storage.NewJSONFileStorage(storeDir)
			default:
				output.EmitError("invalid_store", fmt.Sprintf("Unknown store type: %s", storeType), format)
				return aiverr.InvalidInput("douyin", fmt.Sprintf("unknown store type: %s", storeType))
			}
			if err != nil {
				output.EmitError("storage_error", fmt.Sprintf("Failed to create storage: %v", err), format)
				return err
			}
			defer store.Close()

			// Parse types
			typeList := strings.Split(types, ",")

			// Create pipeline
			p := pipeline.New("douyin", &douyinCollector{client: client}, store)
			if err := p.CollectAndStore(typeList); err != nil {
				output.EmitError("collect_error", fmt.Sprintf("Failed to collect: %v", err), format)
				return err
			}

			return output.EmitSuccess(map[string]interface{}{
				"message": fmt.Sprintf("Collected %d types", len(typeList)),
				"types":   typeList,
			}, format)
		},
	}

	cmd.Flags().StringVar(&types, "types", "hot", "Comma-separated list of types to collect")
	cmd.Flags().StringVar(&storeType, "store", "json", "Storage backend: json or sqlite")
	return cmd
}
