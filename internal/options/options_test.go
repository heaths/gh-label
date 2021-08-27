package options

import "testing"

func Test_RepoOverride(t *testing.T) {
	opts := RootOptions{
		owner: "heaths",
		repo:  "gh-label",
	}

	if owner, repo := opts.RepoOverride(); owner != "heaths" || repo != "gh-label" {
		t.Errorf(`RepoOverride() = (%s, %s); want: ("heaths", "gh-label")`, owner, repo)
	}
}

func Test_parseRepoOverride(t *testing.T) {
	type args struct {
		args string
		env  map[string]string
	}

	type want struct {
		owner string
		repo  string
	}

	tests := []struct {
		name  string
		args  args
		want  want
		wantE bool
	}{
		{
			name: "empty",
			want: want{
				owner: ":owner",
				repo:  ":repo",
			},
		},
		{
			name: "from environment",
			args: args{
				env: map[string]string{
					"GH_REPO": "heaths/gh-label",
				},
			},
			want: want{
				owner: "heaths",
				repo:  "gh-label",
			},
		},
		{
			name: "too few slashes",
			args: args{
				args: "heaths",
			},
			wantE: true,
		},
		{
			name: "too many slashes",
			args: args{
				args: "github.com/heaths/gh-label",
			},
			wantE: true,
		},
		{
			name: "empty parts",
			args: args{
				args: "/",
			},
			wantE: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := RootOptions{
				env: &mockConfig{
					env: tt.args.env,
				},
			}

			if err := opts.parseRepoOverride(tt.args.args); (err != nil) != tt.wantE {
				t.Errorf("parseRepoOverride() = %v, wantE %v", err, tt.wantE)
				return
			}

			if opts.owner != tt.want.owner {
				t.Errorf("parseRepoOverride() owner = %q, want %q", opts.owner, tt.want.owner)
				return
			}

			if opts.repo != tt.want.repo {
				t.Errorf("parseRepoOverride() repo = %q, want %q", opts.repo, tt.want.repo)
			}
		})
	}
}

type mockConfig struct {
	env map[string]string
}

func (c *mockConfig) get(key string) string {
	return c.env[key]
}
