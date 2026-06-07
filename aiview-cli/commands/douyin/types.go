package douyin

import (
	"github.com/jackwener/aiview/internal/output"
	"github.com/spf13/cobra"
)

// Client is the interface that the Douyin API client must satisfy for commands.
type Client interface {
	PlatformName() string
	GetHotSearch() (map[string]interface{}, error)
	GetTrending() (map[string]interface{}, error)
	Search(keyword string, page int, count int) (map[string]interface{}, error)
}

// Credential holds Douyin authentication data.
type Credential struct {
	Cookie string `json:"cookie"`
}

// Helper functions for parsing map responses
func getString(m map[string]interface{}, keys ...string) string {
	if m == nil {
		return ""
	}
	val, ok := m[keys[0]]
	if !ok {
		return ""
	}
	if len(keys) == 1 {
		s, _ := val.(string)
		return s
	}
	sub, ok := val.(map[string]interface{})
	if !ok {
		return ""
	}
	return getString(sub, keys[1:]...)
}

func getInt(m map[string]interface{}, keys ...string) int {
	if m == nil {
		return 0
	}
	val, ok := m[keys[0]]
	if !ok {
		return 0
	}
	if len(keys) == 1 {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
		return 0
	}
	sub, ok := val.(map[string]interface{})
	if !ok {
		return 0
	}
	return getInt(sub, keys[1:]...)
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := val.(map[string]interface{})
	return sub
}

func getSlice(m map[string]interface{}, key string) []interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok {
		return nil
	}
	sub, _ := val.([]interface{})
	return sub
}

// GetOutputFormat extracts the output format from cobra command flags.
func GetOutputFormat(cmd *cobra.Command) output.Format {
	parent := cmd
	for parent.HasParent() {
		parent = parent.Parent()
	}
	asJSON, _ := parent.Flags().GetBool("json")
	asYAML, _ := parent.Flags().GetBool("yaml")
	return output.ResolveFormat(asJSON, asYAML)
}