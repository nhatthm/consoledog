package consoledog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestingT_LastError(t *testing.T) {
	t.Parallel()

	newT := &testingT{}
	newT.Errorf("error: %s", "unknown")

	assert.EqualError(t, newT.LastError(), `error: unknown`)
}
