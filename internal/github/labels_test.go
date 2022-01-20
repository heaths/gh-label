package github

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/MakeNowJust/heredoc"
)

func TestSupportedOutputFormats(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		wantE bool
	}{
		{
			name: "csv",
			want: "csv",
		},
		{
			name: "CSV",
			want: "csv",
		},
		{
			name: ".csv",
			want: "csv",
		},
		{
			name: "json",
			want: "json",
		},
		{
			name: "JSON",
			want: "json",
		},
		{
			name: ".json",
			want: "json",
		},
		{
			name:  "unknown",
			wantE: true,
		},
		{
			name:  "",
			wantE: true,
		},
		{
			name:  "a",
			wantE: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SupportedOutputFormat(tt.name); (err != nil) != tt.wantE {
				t.Errorf("SupportedOutputFormat() error = %v, expected error %v", err, tt.wantE)
			} else if got != tt.want {
				t.Errorf("SupportedOutputFormat() = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormats(t *testing.T) {
	got := OutputFormats()
	want := []string{"csv", "json"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OutputFormats() = %v, expected %v", got, want)
	}
}

func TestLabel_strings(t *testing.T) {
	label := Label{
		Name:        "test",
		Color:       "FF0000",
		Description: "a test",
		URL:         "https://github.com",
	}

	got := label.strings()
	want := []string{
		"test",
		"FF0000",
		"a test",
		"https://github.com",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Strings() = %v, expected %v", got, want)
	}
}

func TestLabels_strings(t *testing.T) {
	labels := Labels{
		Label{
			"foo",
			"FF0000",
			"a foo",
			"https://github.com",
		},
		Label{
			"bar",
			"00FF00",
			"a bar",
			"",
		},
	}

	got := labels.strings()
	want := [][]string{
		{
			"foo",
			"FF0000",
			"a foo",
			"https://github.com",
		},
		{
			"bar",
			"00FF00",
			"a bar",
			"",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Strings() = %v, expected %v", got, want)
	}
}

func TestLabels_write(t *testing.T) {
	labels := Labels{
		Label{
			"foo",
			"FF0000",
			"a foo",
			"https://github.com",
		},
		Label{
			"bar",
			"00FF00",
			"a bar",
			"",
		},
	}

	tests := []struct {
		name   string
		format OutputFormat
		want   string
		wantE  bool
	}{
		{
			name:   "csv",
			format: CSV,
			want: heredoc.Doc(`name,color,description,url
			foo,FF0000,a foo,https://github.com
			bar,00FF00,a bar,
			`),
		},
		{
			name:   "json",
			format: JSON,
			want: heredoc.Doc(`[
			  {
			    "name": "foo",
			    "color": "FF0000",
			    "description": "a foo",
			    "url": "https://github.com"
			  },
			  {
			    "name": "bar",
			    "color": "00FF00",
			    "description": "a bar"
			  }
			]
			`),
		},
		{
			name:   "unknown",
			format: "unknown",
			wantE:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := &bytes.Buffer{}
			if err := labels.Write(tt.format, bytes); (err != nil) != tt.wantE {
				t.Errorf("Write() error = %v, expected error %v", err, tt.wantE)
			} else if bytes.String() != tt.want {
				t.Errorf("Write() = %v, expected %v", bytes.String(), tt.want)
			}
		})
	}
}

func TestReadLabel(t *testing.T) {
	tests := []struct {
		name  string
		data  []string
		want  Label
		wantE bool
	}{
		{
			name:  "too few fields",
			data:  []string{"1", "2", "3"},
			wantE: true,
		},
		{
			name:  "too many fields",
			data:  []string{"1", "2", "3", "4", "right out"},
			wantE: true,
		},
		{
			name: "just right",
			data: []string{"1", "2", "3", "4"},
			want: Label{
				Name:        "1",
				Color:       "2",
				Description: "3",
				URL:         "4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := readLabel(tt.data); (err != nil) != tt.wantE {
				t.Errorf("readLabel() error = %v, expected %v", err, tt.wantE)
			} else if got != nil && !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("readLabel() = %v, expected %v", *got, tt.want)
			}
		})
	}
}
