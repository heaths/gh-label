package gh

import (
	"fmt"
	"testing"
)

func Test_Run(t *testing.T) {
	gh := &Gh{
		runner: &mockGh{
			stdout: "stdout",
			stderr: "stderr",
			err:    fmt.Errorf("error"),
		},
	}

	if stdout, stderr, err := gh.Run(); stdout.String() != "stdout" || stderr.String() != "stderr" || err.Error() != "error" {
		t.Errorf(`Run() = (%q, %q, %q), want ("stdout", "stderr", "error")`, stdout.String(), stderr.String(), err.Error())
	}
}
