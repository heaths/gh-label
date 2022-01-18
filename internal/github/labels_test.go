package github

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/MakeNowJust/heredoc"
)

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
		t.Errorf("Strings() = %v, want: %v", got, want)
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
		t.Errorf("Strings() = %v, want: %v", got, want)
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
			format: 99,
			wantE:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := &bytes.Buffer{}
			if err := labels.Write(tt.format, bytes); (err != nil) != tt.wantE {
				t.Errorf("Write() error = %v, expected error", err)
			} else if bytes.String() != tt.want {
				t.Errorf("Write() = %v, want: %v", bytes.String(), tt.want)
			}
		})
	}
}
