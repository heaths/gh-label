package create

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
)

func Test_create(t *testing.T) {
	tests := []struct {
		name string
		tty  bool
		want string
	}{
		{
			name: "create",
			want: "https://github.com/heaths/gh-label/labels/test\n",
		},
		{
			name: "create (TTY)",
			tty:  true,
			want: heredoc.Doc(`Created label
			
			https://github.com/heaths/gh-label/labels/test
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
			mock := &github.Mock{
				Stdout: *bytes.NewBufferString(heredoc.Doc(`{
					"id": 3315930645,
					"node_id": "MDU6TGFiZWwzMzE1OTMwNjQ1",
					"url": "https://api.github.com/repos/heaths/gh-label/labels/test",
					"name": "test",
					"color": "112233",
					"default": false,
					"description": ""
			}`)),
			}

			rootOpts := &options.GlobalOptions{}
			opts := &createOptions{
				name:  "test",
				color: "112233",

				client: github.New(mock),
				io:     io,
			}

			if err := create(rootOpts, opts); err != nil {
				t.Errorf("create() error = %v", err)
				return
			}

			if gotW := stdout.String(); gotW != tt.want {
				t.Errorf("create() = %q, want %q", gotW, tt.want)
			}
		})
	}
}
