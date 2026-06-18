package bilibili

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/output"
	"github.com/jackwener/aiview/internal/pipeline"
	"github.com/jackwener/aiview/internal/storage"
	"github.com/spf13/cobra"
)

// bilibiliCollector wraps the Client to implement pipeline.Collector
type bilibiliCollector struct {
	client Client
}

func (bc *bilibiliCollector) Collect(recordType string) (map[string]interface{}, error) {
	switch recordType {
	case "hot":
		videos, err := bc.client.GetHotVideos(1, 20)
		if err != nil {
			return nil, err
		}
		return videosToMap(videos), nil
	case "trending":
		videos, err := bc.client.GetRankVideos(0, 3, "all")
		if err != nil {
			return nil, err
		}
		return videosToMap(videos), nil
	case "search":
		results, err := bc.client.SearchVideo("热门", 1, "totalrank", 0, 0)
		if err != nil {
			return nil, err
		}
		return searchResultsToMap(results), nil
	default:
		return nil, aiverr.InvalidInput("bilibili", fmt.Sprintf("unknown type: %s", recordType))
	}
}

func videosToMap(videos interface{}) map[string]interface{} {
	data, _ := json.Marshal(videos)
	var result []interface{}
	json.Unmarshal(data, &result)
	return map[string]interface{}{"items": result}
}

func searchResultsToMap(results interface{}) map[string]interface{} {
	data, _ := json.Marshal(results)
	var result []interface{}
	json.Unmarshal(data, &result)
	return map[string]interface{}{"items": result}
}

// NewCollectCmd creates the collect command for bilibili.
func NewCollectCmd(getClient func() Client) *cobra.Command {
	var types string
	var storeType string

	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Batch collect and store data",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getClient()
			format := output.MustGetFormat(cmd)

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
				return aiverr.InvalidInput("bilibili", fmt.Sprintf("unknown store type: %s", storeType))
			}
			if err != nil {
				output.EmitError("storage_error", fmt.Sprintf("Failed to create storage: %v", err), format)
				return err
			}
			defer store.Close()

			// Parse types
			typeList := strings.Split(types, ",")

			// Create pipeline
			p := pipeline.New("bilibili", &bilibiliCollector{client: client}, store)
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
