package consoledog

import (
	"fmt"
	"strings"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
)

type testingT struct {
	err error
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.err = fmt.Errorf(format, args...) // nolint: goerr113
}

func (t *testingT) LastError() error {
	return t.err
}

func t() *testingT {
	return &testingT{}
}

// AssertState asserts console state.
func AssertState(t assert.TestingT, state *vt10x.State, expected string) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	actual := trimTailingSpaces(expect.StripTrailingEmptyLines(state.String()))

	return assert.Equal(t, expected, actual)
}

func trimTailingSpaces(out string) string {
	lines := strings.Split(out, "\n")

	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}
