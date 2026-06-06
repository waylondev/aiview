package cli

import (
	"fmt"
	"os"

	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/output"
	"github.com/jackwener/aiview/internal/platform"
	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	asJSON    bool
	asYAML    bool
	verbose   bool
)

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "aiview",
	Short: "aiview — 多平台 CLI 工具",
	Long:  `aiview 是一个多平台 CLI 工具，支持 Bilibili 等平台的内容浏览、搜索和交互。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	// Register all platform commands (must happen after all init() have run)
	for _, p := range platform.All() {
		for _, cmd := range p.Commands() {
			rootCmd.AddCommand(cmd)
		}
	}
	return rootCmd.Execute()
}

// GetConfig returns the loaded configuration.
func GetConfig() *config.Config {
	return cfg
}

// GetOutputFormat resolves the current output format.
func GetOutputFormat() output.Format {
	return output.ResolveFormat(asJSON, asYAML)
}

// IsVerbose returns whether verbose logging is enabled.
func IsVerbose() bool {
	return verbose
}

// ExitError prints an error message and exits with code 1.
func ExitError(code, message string) {
	format := output.ResolveFormat(asJSON, asYAML)
	output.EmitError(code, message, format)
	os.Exit(1)
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&asJSON, "json", false, "输出 JSON 格式")
	rootCmd.PersistentFlags().BoolVar(&asYAML, "yaml", false, "输出 YAML 格式")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "启用详细日志")
}