package github

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/heaths/gh-label/internal/utils"
)

const labelFields = 4

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

func ReadLabels(format OutputFormat, r io.Reader) (Labels, error) {
	// Start with capacity for 10 labels. A new repo currently starts with 9.
	labels := make(Labels, 0, 10)

	if format == CSV {
		csv := csv.NewReader(r)
		csv.FieldsPerRecord = labelFields
		csv.ReuseRecord = true
		csv.TrimLeadingSpace = true

		for {
			record, err := csv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			if utils.AreEqualStrings(record, labels.headers()) {
				continue
			}

			label, err := readLabel(record)
			if err != nil {
				return nil, err
			}

			labels = append(labels, *label)
		}

		return labels, nil
	}

	if format == JSON {
		json := json.NewDecoder(r)
		if err := json.Decode(&labels); err != nil {
			return nil, err
		}

		return labels, nil
	}

	return nil, fmt.Errorf("unknown format %v", format)
}

func readLabel(record []string) (*Label, error) {
	if len(record) != labelFields {
		return nil, fmt.Errorf("expected %d label fields, got %d", labelFields, len(record))
	}

	label := &Label{
		record[0],
		record[1],
		record[2],
		record[3],
	}

	return label, nil
}
