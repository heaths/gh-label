package edit

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
)

func Test_edit(t *testing.T) {
	tests := []struct {
		name   string
		rename string
		tty    bool
		want   string
	}{
		{
			name: "edit",
			want: "https://github.com/heaths/gh-label/labels/test\n",
		},
		{
			name: "edit (TTY)",
			tty:  true,
			want: heredoc.Doc(`Updated label 'test'
			
			https://github.com/heaths/gh-label/labels/test
			`),
		},
		{
			name:   "renamed (TTY)",
			rename: "test2",
			tty:    true,
			want: heredoc.Doc(`Renamed label 'test' to 'test2'
			
			https://github.com/heaths/gh-label/labels/test2
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up output streams.
			io, _, stdout, _ := iostreams.Test()
			io.SetStdoutTTY(tt.tty)
			io.SetColorEnabled(true)

			// Set up gh output.
			name := "test"
			if tt.rename != "" {
				name = tt.rename
			}
			mock := &github.Mock{
				Stdout: *bytes.NewBufferString(heredoc.Docf(`{
					"id": 3315930645,
					"node_id": "MDU6TGFiZWwzMzE1OTMwNjQ1",
					"url": "https://api.github.com/repos/heaths/gh-label/labels/%[1]s",
					"name": "%[1]s",
					"color": "112233",
					"default": false,
					"description": ""
				}`, name)),
			}

			rootOpts := &options.GlobalOptions{}
			opts := &editOptions{
				name:    "test",
				color:   "112233",
				newName: tt.rename,

				client: github.New(mock),
				io:     io,
			}

			if err := edit(rootOpts, opts); err != nil {
				t.Errorf("edit() error = %v", err)
				return
			}

			if gotW := stdout.String(); gotW != tt.want {
				t.Errorf("edit() = %q, want %q", gotW, tt.want)
			}
		})
	}
}
