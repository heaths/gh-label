package github

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Labels []Label

type OutputFormat string

const (
	CSV  OutputFormat = "csv"
	JSON OutputFormat = "json"
)

func SupportedOutputFormat(format string) (string, error) {
	format = strings.TrimPrefix(format, ".")
	format = strings.ToLower(format)
	formats := OutputFormats()

	for _, str := range formats {
		if str == format {
			return format, nil
		}
	}

	return "", fmt.Errorf("unsupported format %q, expected %v", format, formats)
}

func OutputFormats() []string {
	// These must remain sorted.
	return []string{"csv", "json"}
}

func (label *Label) strings() []string {
	return []string{
		label.Name,
		label.Color,
		label.Description,
		label.URL,
	}
}

func (labels *Labels) headers() []string {
	return []string{
		"name",
		"color",
		"description",
		"url",
	}
}

func (labels *Labels) strings() [][]string {
	arr := make([][]string, len(*labels))
	for i, elem := range *labels {
		arr[i] = elem.strings()
	}
	return arr
}

func (labels *Labels) Write(format OutputFormat, w io.Writer) error {
	if format == CSV {
		csv := csv.NewWriter(w)
		if err := csv.Write(labels.headers()); err != nil {
			return err
		}
		return csv.WriteAll(labels.strings())
	}
	if format == JSON {
		json := json.NewEncoder(w)
		json.SetIndent("", "  ")
		return json.Encode(*labels)
	}
	return fmt.Errorf("unknown format %v", format)
}
