package list

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/gh"
	"github.com/heaths/gh-label/internal/options"
)

func Test_list(t *testing.T) {
	type args struct {
		stdout string
		stderr string
		tty    bool
	}

	tests := []struct {
		name  string
		args  args
		wantW string
		wantE bool
	}{
		{
			name: "single page",
			args: args{
				stdout: `{
					"data": {
						"repository": {
							"labels": {
								"nodes": [
									{
										"name": "bug",
										"color": "d73a4a",
										"description": "Something isn't working"
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
			wantW: heredoc.Docf(`bug%[1]sd73a4a%[1]sSomething isn't working
			documentation%[1]s0075ca%[1]sImprovements or additions to documentation
			`, "\t"),
		},
		{
			name: "single page (TTY)",
			args: args{
				stdout: `{
					"data": {
						"repository": {
							"labels": {
								"nodes": [
									{
										"name": "bug",
										"color": "d73a4a",
										"description": "Something isn't working"
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
				tty: true,
			},
			wantW: heredoc.Docf(`bug            #d73a4a  %[1]s[0;90mSomething isn't working%[1]s[0m
			documentation  #0075ca  %[1]s[0;90mImprovements or additions to documentation%[1]s[0m
			`, "\x1b"),
		},
		{
			name: "multiple pages",
			args: args{
				// Tests workaround for invalid JSON until https://github.com/cli/cli/issues/1268 is resolved.
				stdout: `{
					"data": {
						"repository": {
							"labels": {
								"nodes": [
									{
										"name": "bug",
										"color": "d73a4a",
										"description": "Something isn't working"
									},
									{
										"name": "documentation",
										"color": "0075ca",
										"description": "Improvements or additions to documentation"
									}
								],
								"pageInfo":{"hasNextPage":true,"endCursor":"abcd1234"}}}}}
				{
					"data": {
						"repository": {
							"labels": {
								"nodes": [
									{
										"name": "duplicate",
										"color": "cfd3d7",
										"description": "This issue or pull request already exists"
									},
									{
										"name": "enhancement",
										"color": "a2eeef",
										"description": "New feature or request"
									}
								],
								"pageInfo":{"hasNextPage":false}}}}}`,
			},
			wantW: heredoc.Docf(`bug%[1]sd73a4a%[1]sSomething isn't working
			documentation%[1]s0075ca%[1]sImprovements or additions to documentation
			duplicate%[1]scfd3d7%[1]sThis issue or pull request already exists
			enhancement%[1]sa2eeef%[1]sNew feature or request
			`, "\t"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up output streams.
			io, _, stdout, _ := iostreams.Test()
			io.SetStdoutTTY(tt.args.tty)
			io.SetStderrTTY(tt.args.tty)
			io.SetColorEnabled(true)

			// Set up gh output.
			gh := gh.NewMock(tt.args.stdout, tt.args.stderr, nil)

			rootOpts := &options.RootOptions{}
			opts := &listOptions{
				io: io,
				gh: gh,
			}

			if err := list(rootOpts, opts); (err != nil) != tt.wantE {
				t.Errorf("list() error = %v, wantE %v", err, tt.wantE)
				return
			}

			if gotW := stdout.String(); gotW != tt.wantW {
				t.Errorf("list() = %q, want %q", gotW, tt.wantW)
			}
		})
	}
}
