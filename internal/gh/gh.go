package gh

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/cli/safeexec"
)

type runner interface {
	run(args ...string) (stdout, stderr bytes.Buffer, err error)
}

type Gh struct {
	runner runner
}

func (gh *Gh) Run(args ...string) (stdout, stderr bytes.Buffer, err error) {
	if gh.runner == nil {
		gh.runner = &execRunner{}
	}

	return gh.runner.run(args...)
}

type execRunner struct{}

func (r *execRunner) run(args ...string) (stdout, stderr bytes.Buffer, err error) {
	bin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("cannot find gh; is it installed? err: %w", err)
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh; error: %w, stderr: %s", err, stderr.String())
		return
	}

	return
}

type mockGh struct {
	stdout string
	stderr string
	err    error
}

func NewMock(stdout, stderr string, err error) *Gh {
	return &Gh{
		runner: &mockGh{
			stdout: stdout,
			stderr: stderr,
			err:    err,
		},
	}
}

func (r *mockGh) run(args ...string) (stdout, stderr bytes.Buffer, err error) {
	return *bytes.NewBufferString(r.stdout), *bytes.NewBufferString(r.stderr), r.err
}
