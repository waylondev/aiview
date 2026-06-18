// Package output provides unified output formatting in JSON, YAML, and text modes.
package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// SchemaVersion is the version of the output schema format.
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
	FormatCSV
)

// ResolveFormat determines the output format based on flags and TTY status.
func ResolveFormat(asJSON, asYAML, asTable, asCSV bool) Format {
	count := 0
	for _, v := range []bool{asJSON, asYAML, asTable, asCSV} {
		if v {
			count++
		}
	}
	if count > 1 {
		fmt.Fprintln(os.Stderr, "Cannot use multiple output format flags simultaneously")
		os.Exit(1)
	}
	if asJSON {
		return FormatJSON
	}
	if asYAML {
		return FormatYAML
	}
	if asTable {
		return FormatTable
	}
	if asCSV {
		return FormatCSV
	}
	outputMode := strings.ToLower(os.Getenv("OUTPUT"))
	switch outputMode {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	case "csv":
		return FormatCSV
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

// EmitPlatformError outputs a structured PlatformError with full details.
func EmitPlatformError(err error, format Format) bool {
	if pe, ok := aiverr.IsPlatformError(err); ok {
		env := Envelope{
			OK:            false,
			SchemaVersion: SchemaVersion,
			Error: &ErrorInfo{
				Code:    pe.Code,
				Message: pe.Message,
				Details: pe.Details,
			},
		}
		emitEnvelope(env, format)
		return false
	}
	// Fallback for non-PlatformError
	return EmitError("unknown_error", err.Error(), format)
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
	case FormatCSV:
		// For CSV format, commands should handle their own CSV rendering
		return nil
	}
	return nil
}

// NewTableWriter creates a new tabwriter for aligned table output.
func NewTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
}

// RenderTable renders data as a table using tabwriter.
// data should be a slice of structs or a slice of slices.
func RenderTable(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return aiverr.InvalidInput("output", fmt.Sprintf("RenderTable: expected slice, got %T", data))
	}
	if v.Len() == 0 {
		return nil
	}

	w := NewTableWriter()
	defer w.Flush()

	// Check if it's a slice of slices
	if v.Index(0).Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			row := v.Index(i)
			for j := 0; j < row.Len(); j++ {
				fmt.Fprintf(w, "%v", row.Index(j).Interface())
				if j < row.Len()-1 {
					fmt.Fprint(w, "\t")
				}
			}
			fmt.Fprintln(w)
		}
		return nil
	}

	// It's a slice of structs - extract headers from first element
	firstElem := v.Index(0)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}
	if firstElem.Kind() != reflect.Struct {
		return aiverr.InvalidInput("output", fmt.Sprintf("RenderTable: expected slice of structs or slices, got %T", data))
	}

	t := firstElem.Type()
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			headers = append(headers, field.Name)
		}
	}

	// Print headers
	for i, h := range headers {
		fmt.Fprint(w, h)
		if i < len(headers)-1 {
			fmt.Fprint(w, "\t")
		}
	}
	fmt.Fprintln(w)

	// Print rows
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		for j, h := range headers {
			field := elem.FieldByName(h)
			fmt.Fprintf(w, "%v", field.Interface())
			if j < len(headers)-1 {
				fmt.Fprint(w, "\t")
			}
		}
		fmt.Fprintln(w)
	}

	return nil
}

// RenderCSV renders data as CSV using encoding/csv.
// data should be a slice of structs or a slice of slices.
func RenderCSV(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return aiverr.InvalidInput("output", fmt.Sprintf("RenderCSV: expected slice, got %T", data))
	}
	if v.Len() == 0 {
		return nil
	}

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// Check if it's a slice of slices
	if v.Index(0).Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			row := v.Index(i)
			var record []string
			for j := 0; j < row.Len(); j++ {
				record = append(record, fmt.Sprintf("%v", row.Index(j).Interface()))
			}
			if err := w.Write(record); err != nil {
				return err
			}
		}
		return nil
	}

	// It's a slice of structs - extract headers from first element
	firstElem := v.Index(0)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}
	if firstElem.Kind() != reflect.Struct {
		return aiverr.InvalidInput("output", fmt.Sprintf("RenderCSV: expected slice of structs or slices, got %T", data))
	}

	t := firstElem.Type()
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			headers = append(headers, field.Name)
		}
	}

	// Write headers
	if err := w.Write(headers); err != nil {
		return err
	}

	// Write rows
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		var record []string
		for _, h := range headers {
			field := elem.FieldByName(h)
			record = append(record, fmt.Sprintf("%v", field.Interface()))
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// FormatCount formats a large number with a "w" (万) suffix.
func FormatCount(n int) string {
	if n >= 10000 {
		return fmt.Sprintf("%.1f万", float64(n)/10000)
	}
	return fmt.Sprintf("%d", n)
}

// GetFormat extracts the output format from cobra command flags.
func GetFormat(cmd *cobra.Command) Format {
	parent := cmd
	for parent.HasParent() {
		parent = parent.Parent()
	}
	asJSON, _ := parent.Flags().GetBool("json")
	asYAML, _ := parent.Flags().GetBool("yaml")
	asTable, _ := parent.Flags().GetBool("table")
	asCSV, _ := parent.Flags().GetBool("csv")
	return ResolveFormat(asJSON, asYAML, asTable, asCSV)
}