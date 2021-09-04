package delete

import (
	"errors"
	"testing"

	"github.com/cli/cli/pkg/iostreams"
	"github.com/heaths/gh-label/internal/github"
	"github.com/heaths/gh-label/internal/options"
)

func Test_delete(t *testing.T) {
	t.Run("delete (TTY)", func(t *testing.T) {
		// Set up output streams.
		io, _, stdout, _ := iostreams.Test()
		io.SetStdoutTTY(true)
		io.SetColorEnabled(true)

		// Set up gh output.
		mock := &github.Mock{}

		rootOpts := &options.GlobalOptions{}
		opts := &deleteOptions{
			name: "test",

			client: github.New(mock),
			io:     io,
		}

		if err := delete(rootOpts, opts); err != nil {
			t.Errorf("create() error = %v", err)
			return
		}

		want := "Deleted label 'test'\n"
		if want != stdout.String() {
			t.Errorf("create() = %s, want: %s", stdout.String(), want)
		}
	})
}

func Test_delete_error(t *testing.T) {
	t.Run("delete with error", func(t *testing.T) {
		// Set up output streams.
		io, _, _, _ := iostreams.Test()
		io.SetStdoutTTY(true)
		io.SetColorEnabled(true)

		// Set up gh output.
		mock := &github.Mock{
			Err: errors.New("gh returned error: exit status 1, stderr: gh: Not Found (HTTP 404)"),
		}

		rootOpts := &options.GlobalOptions{}
		opts := &deleteOptions{
			name: "test",

			client: github.New(mock),
			io:     io,
		}

		if err := delete(rootOpts, opts); err == nil {
			t.Errorf("create() error = nil, expected error")
			return
		}
	})
}
