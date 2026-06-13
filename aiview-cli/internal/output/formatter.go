// Package output provides unified output formatting in JSON, YAML, and text modes.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

const SchemaVersion = "1"

// Envelope is the agent-friendly output envelope.
type Envelope struct {
	OK            bool        `json:"ok" yaml:"ok"`
	SchemaVersion string      `json:"schema_version" yaml:"schema_version"`
	Data          interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	Error         *ErrorInfo  `json:"error,omitempty" yaml:"error,omitempty"`
}

// ErrorInfo holds structured error information.
type ErrorInfo struct {
	Code    string      `json:"code" yaml:"code"`
	Message string      `json:"message" yaml:"message"`
	Details interface{} `json:"details,omitempty" yaml:"details,omitempty"`
}

// Format represents the output format.
type Format int

const (
	FormatAuto Format = iota
	FormatJSON
	FormatYAML
	FormatTable
)

// ResolveFormat determines the output format based on flags and TTY status.
func ResolveFormat(asJSON, asYAML bool) Format {
	if asJSON && asYAML {
		fmt.Fprintln(os.Stderr, "Cannot use both --json and --yaml")
		os.Exit(1)
	}
	if asJSON {
		return FormatJSON
	}
	if asYAML {
		return FormatYAML
	}
	outputMode := strings.ToLower(os.Getenv("OUTPUT"))
	switch outputMode {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	case "rich", "table":
		return FormatTable
	}
	if !isTTY() {
		return FormatYAML
	}
	return FormatTable
}

// isTTY checks if stdout is a terminal.
func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// EmitSuccess outputs a success envelope in the specified format.
func EmitSuccess(data interface{}, format Format) error {
	env := Envelope{
		OK:            true,
		SchemaVersion: SchemaVersion,
		Data:          data,
	}
	return emitEnvelope(env, format)
}

// EmitError outputs an error envelope and returns false (for chaining).
func EmitError(code, message string, format Format) bool {
	env := Envelope{
		OK:            false,
		SchemaVersion: SchemaVersion,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
	emitEnvelope(env, format)
	return false
}

func emitEnvelope(env Envelope, format Format) error {
	switch format {
	case FormatJSON:
		data, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case FormatYAML:
		data, err := yaml.Marshal(env)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	case FormatTable:
		// For table format, we just print the data directly
		// Commands should handle their own table rendering
		return nil
	}
	return nil
}

// NewTableWriter creates a new tabwriter for aligned table output.
func NewTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
}

// FormatCount formats a large number with a "w" (万) suffix.
func FormatCount(n int) string {
	if n >= 10000 {
		return fmt.Sprintf("%.1f万", float64(n)/10000)
	}
	return fmt.Sprintf("%d", n)
}

// FormatDuration formats seconds into MM:SS or HH:MM:SS.
func FormatDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	if seconds >= 3600 {
		h := seconds / 3600
		m := (seconds % 3600) / 60
		s := seconds % 60
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}