package importcmd

// cSpell:ignore fstest

import (
	"bytes"
	"errors"
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
)

var (
	csvData = []byte(heredoc.Doc(`name,color,description,url
		bug,d73a4a,Something isn't working,https://github.com/heaths/gh-label/issues/1
		`))

	jsonData = []byte(heredoc.Doc(`[
			{
				"name": "bug",
				"color": "d73a4a",
				"description": "Something isn't working",
				"url": "https://github.com/heaths/gh-label/issues/1"
			}
		]`))

	jsonLabel = bytes.Trim(jsonData, "[]")
)

func Test_ImportCmd(t *testing.T) {
	type args struct {
		path   string
		format string
		data   []byte
	}

	tests := []struct {
		name  string
		args  args
		want  args
		wantE bool
	}{
		{
			name: "csv file",
			args: args{
				path: "labels.csv",
				data: csvData,
			},
			want: args{
				path:   "labels.csv",
				format: "csv",
			},
		},
		{
			name: "csv stream",
			args: args{
				path:   "-",
				format: "csv",
				data:   csvData,
			},
			want: args{
				path:   "-",
				format: "csv",
			},
		},
		{
			name: "csv override",
			args: args{
				path:   "labels.json",
				format: "csv",
				data:   csvData,
			},
			want: args{
				path:   "labels.json",
				format: "csv",
			},
		},
		{
			name: "json file",
			args: args{
				path: "labels.json",
				data: jsonData,
			},
			want: args{
				path:   "labels.json",
				format: "json",
			},
		},
		{
			name: "json stream",
			args: args{
				path:   "-",
				format: "json",
				data:   jsonData,
			},
			want: args{
				path:   "-",
				format: "json",
			},
		},
		{
			name: "json override",
			args: args{
				path:   "labels.csv",
				format: "json",
				data:   jsonData,
			},
			want: args{
				path:   "labels.csv",
				format: "json",
			},
		},
		{
			name: "stream without format",
			args: args{
				path: "-",
			},
			want: args{
				path: "-",
			},
			wantE: true,
		},
		{
			name: "invalid format",
			args: args{
				path:   "-",
				format: "",
			},
			want: args{
				path: "-",
			},
			wantE: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globalOpts := &options.GlobalOptions{}
			cmd := ImportCmd(globalOpts)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			args := []string{tt.args.path}
			if tt.args.format != "" {
				args = append(args, "--format", tt.args.format)
			}
			cmd.SetArgs(args)

			mock := &github.Mock{
				Stdout: *bytes.NewBuffer(jsonLabel),
			}
			opts.client = github.New(mock)

			fs := fstest.MapFS{}
			fs[tt.args.path] = &fstest.MapFile{
				Data: tt.args.data,
			}
			opts.fs = fs

			io, stdin, _, _ := iostreams.Test()
			opts.io = io
			if tt.args.path == "-" {
				stdin.Write(tt.args.data)
			}

			if err := cmd.Execute(); (err != nil) != tt.wantE {
				t.Errorf("ImportCmd().Execute() error = %v, expected %v", err, tt.wantE)
				return
			}

			if opts.path != tt.want.path {
				t.Errorf("ImportCmd() path = %q, expected %q", opts.path, tt.want.path)
				return
			}

			if opts.format != tt.want.format {
				t.Errorf("ImportCmd() format = %q, expected %q", opts.format, tt.want.format)
				return
			}
		})
	}
}

func Test_import(t *testing.T) {
	type args struct {
		format string
		stdin  []byte
		tty    bool
		err    error
	}

	tests := []struct {
		name  string
		args  args
		wantW string
		wantE bool
	}{
		{
			name: "csv",
			args: args{
				format: "csv",
				stdin:  csvData,
			},
		},
		{
			name: "csv (tty)",
			args: args{
				format: "csv",
				stdin:  csvData,
				tty:    true,
			},
			wantW: heredoc.Doc(`Importing 1 label(s) from "-"

			Successfully imported 1, failed to import 0 label(s)
			`),
		},
		{
			name: "json",
			args: args{
				format: "json",
				stdin:  jsonData,
			},
		},
		{
			name: "json (tty)",
			args: args{
				format: "json",
				stdin:  jsonData,
				tty:    true,
			},
			wantW: heredoc.Doc(`Importing 1 label(s) from "-"

			Successfully imported 1, failed to import 0 label(s)
			`),
		},
		{
			name: "all failed",
			args: args{
				format: "json",
				stdin:  jsonData,
				err:    errors.New("failed"),
			},
			wantE: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up streams.
			io, stdin, stdout, _ := iostreams.Test()
			io.SetStdoutTTY(tt.args.tty)
			stdin.Write(tt.args.stdin)

			// Set up gh output.
			mock := &github.Mock{
				Stdout: *bytes.NewBuffer(jsonLabel),
				Err:    tt.args.err,
			}

			rootOpts := &options.GlobalOptions{}
			opts := &importOptions{
				path:   "-",
				format: tt.args.format,

				client: github.New(mock),
				io:     io,
			}

			if err := _import(rootOpts, opts); (err != nil) != tt.wantE {
				t.Errorf("_import() error = %v, wantE %v", err, tt.wantE)
				return
			}

			if gotW := stdout.String(); gotW != tt.wantW {
				t.Errorf("_import() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}
