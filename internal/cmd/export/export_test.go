package export

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
)

func Test_ExportCmd(t *testing.T) {
	type args struct {
		path   string
		format string
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
			cmd := ExportCmd(globalOpts)
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			args := []string{tt.args.path}
			if tt.args.format != "" {
				args = append(args, "--format", tt.args.format)
			}
			cmd.SetArgs(args)

			mock := &github.Mock{
				Err: fmt.Errorf("mock"),
			}
			opts.client = github.New(mock)

			if err := cmd.Execute(); !strings.Contains(err.Error(), "mock") && (err != nil) != tt.wantE {
				t.Error("ImportCmd().Execute() expected error")
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

func Test_export(t *testing.T) {
	type args struct {
		format string
		stdout string
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
				stdout: `{
					"data": {
						"repository": {
							"labels": {
								"nodes": [
									{
										"name": "bug",
										"color": "d73a4a",
										"description": "Something isn't working",
										"url": "https://github.com/heaths/gh-label/issues/1"
									},
									{
										"name": "documentation",
										"color": "0075ca",
										"description": "Improvements or additions to documentation"
									}
								]
							}
						}
					}
				}`,
			},
			wantW: heredoc.Doc(`name,color,description,url
			bug,d73a4a,Something isn't working,https://github.com/heaths/gh-label/issues/1
			documentation,0075ca,Improvements or additions to documentation,
			`),
		},
		{
			name: "json",
			args: args{
				format: "json",
				stdout: `{
					"data": {
						"repository": {
							"labels": {
								"nodes": [
									{
										"name": "bug",
										"color": "d73a4a",
										"description": "Something isn't working",
										"url": "https://github.com/heaths/gh-label/issues/1"
									},
									{
										"name": "documentation",
										"color": "0075ca",
										"description": "Improvements or additions to documentation"
									}
								]
							}
						}
					}
				}`,
			},
			wantW: heredoc.Doc(`[
			  {
			    "name": "bug",
			    "color": "d73a4a",
			    "description": "Something isn't working",
			    "url": "https://github.com/heaths/gh-label/issues/1"
			  },
			  {
			    "name": "documentation",
			    "color": "0075ca",
			    "description": "Improvements or additions to documentation"
			  }
			]
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up output streams.
			io, _, stdout, _ := iostreams.Test()

			// Set up gh output.
			mock := &github.Mock{
				Stdout: *bytes.NewBufferString(tt.args.stdout),
			}

			rootOpts := &options.GlobalOptions{}
			opts := &exportOptions{
				path:   "-",
				format: tt.args.format,

				client: github.New(mock),
				io:     io,
			}

			if err := export(rootOpts, opts); (err != nil) != tt.wantE {
				t.Errorf("export() error = %v, wantE %v", err, tt.wantE)
				return
			}

			if gotW := stdout.String(); gotW != tt.wantW {
				t.Errorf("export() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}
