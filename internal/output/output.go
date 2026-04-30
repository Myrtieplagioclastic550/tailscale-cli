package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// Print formats and prints data in the specified format.
// Supported formats: table, json, yaml, csv.
// data should be a []map[string]interface{} or a map[string]interface{}.
// columns defines which keys to display and their order.
func Print(format string, data interface{}, columns []string) error {
	switch format {
	case "json":
		return PrintJSON(data)
	case "yaml":
		return printYAML(data)
	case "table":
		return printTable(data, columns)
	case "csv":
		return printCSV(data, columns)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// PrintJSON outputs data as indented JSON.
func PrintJSON(data interface{}) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	fmt.Println(string(out))
	return nil
}

func printYAML(data interface{}) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}
	fmt.Print(string(out))
	return nil
}

// toRows normalizes data into a slice of maps.
func toRows(data interface{}) ([]map[string]interface{}, error) {
	switch v := data.(type) {
	case []map[string]interface{}:
		return v, nil
	case map[string]interface{}:
		return []map[string]interface{}{v}, nil
	case []interface{}:
		rows := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			m, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected map[string]interface{}, got %T", item)
			}
			rows = append(rows, m)
		}
		return rows, nil
	default:
		return nil, fmt.Errorf("unsupported data type for table/csv output: %T", data)
	}
}

func printTable(data interface{}, columns []string) error {
	rows, err := toRows(data)
	if err != nil {
		return err
	}

	if len(columns) == 0 && len(rows) > 0 {
		for k := range rows[0] {
			columns = append(columns, k)
		}
	}

	table := tablewriter.NewTable(os.Stdout)

	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = strings.ToUpper(col)
	}
	table.Header(headers)

	for _, row := range rows {
		vals := make([]string, len(columns))
		for i, col := range columns {
			vals[i] = fmt.Sprintf("%v", row[col])
		}
		table.Append(vals)
	}

	table.Render()
	return nil
}

func printCSV(data interface{}, columns []string) error {
	rows, err := toRows(data)
	if err != nil {
		return err
	}

	if len(columns) == 0 && len(rows) > 0 {
		for k := range rows[0] {
			columns = append(columns, k)
		}
	}

	// Print header
	fmt.Println(strings.Join(columns, ","))

	// Print rows
	for _, row := range rows {
		vals := make([]string, len(columns))
		for i, col := range columns {
			vals[i] = fmt.Sprintf("%v", row[col])
		}
		fmt.Println(strings.Join(vals, ","))
	}

	return nil
}
