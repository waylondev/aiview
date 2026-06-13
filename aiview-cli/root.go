package main

import (
	"fmt"
	"os"

	"github.com/jackwener/aiview/commands"
	"github.com/jackwener/aiview/internal/config"
	"github.com/jackwener/aiview/internal/output"
	"github.com/jackwener/aiview/internal/platform"
	_ "github.com/jackwener/aiview/internal/platform/bilibili"
	_ "github.com/jackwener/aiview/internal/platform/douyin"
	_ "github.com/jackwener/aiview/internal/platform/xiaohongshu"
	"github.com/spf13/cobra"
)

var (
	cfg     *config.Config
	asJSON  bool
	asYAML  bool
	asTable bool
	asCSV   bool
	verbose bool
)

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "aiview",
	Short: "aiview — Multi-platform CLI tool",
	Long:  `aiview is a multi-platform CLI tool for browsing, searching, and interacting with content from Bilibili and other platforms.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	// Register global commands
	rootCmd.AddCommand(commands.NewExportCmd())
	rootCmd.AddCommand(commands.NewAnalyzeCmd())
	rootCmd.AddCommand(commands.NewCompareCmd())
	rootCmd.AddCommand(commands.NewScheduleCmd())

	// Register all platform commands (must happen after all init() have run)
	for _, p := range platform.All() {
		for _, cmd := range p.Commands() {
			rootCmd.AddCommand(cmd)
		}
	}
	return rootCmd.Execute()
}

// getConfig returns the loaded configuration.
func getConfig() *config.Config {
	return cfg
}

// getOutputFormat resolves the current output format.
func getOutputFormat() output.Format {
	return output.ResolveFormat(asJSON, asYAML, asTable, asCSV)
}

// isVerbose returns whether verbose logging is enabled.
func isVerbose() bool {
	return verbose
}

// exitError prints an error message and exits with code 1.
func exitError(code, message string) {
	format := output.ResolveFormat(asJSON, asYAML, asTable, asCSV)
	output.EmitError(code, message, format)
	os.Exit(1)
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&asJSON, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&asYAML, "yaml", false, "Output in YAML format")
	rootCmd.PersistentFlags().BoolVar(&asTable, "table", false, "Output in table format")
	rootCmd.PersistentFlags().BoolVar(&asCSV, "csv", false, "Output in CSV format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}