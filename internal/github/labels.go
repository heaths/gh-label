package github

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Labels []Label

type OutputFormat int

const (
	CSV OutputFormat = iota
	JSON
)

func (label *Label) strings() []string {
	return []string{
		label.Name,
		label.Color,
		label.Description,
		label.URL,
	}
}

func (label *Labels) headers() []string {
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
		csv.Write(labels.headers())
		return csv.WriteAll(labels.strings())
	}
	if format == JSON {
		json := json.NewEncoder(w)
		json.SetIndent("", "  ")
		return json.Encode(*labels)
	}
	return fmt.Errorf("unknown format %v", format)
}
