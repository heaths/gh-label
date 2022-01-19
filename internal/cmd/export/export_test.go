package export

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
)

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
