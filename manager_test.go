package consoledog_test

import (
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/cucumber/godog"
	"github.com/nhatthm/consoledog"
)

func TestManager(t *testing.T) {
	t.Parallel()

	m := consoledog.New(t,
		consoledog.WithTermSize(80, 24),
		consoledog.WithStarter(func(sc *godog.Scenario, console *expect.Console) {
			console.Write([]byte(`hello world`)) // nolint: errcheck, gosec
		}),
	)

	scenario := &godog.Scenario{}
	_, state := m.NewConsole(scenario)

	// New again does not affect the state.
	_, _ = m.NewConsole(scenario)

	m.Flush()

	expected := `hello world`

	consoledog.AssertState(t, state, expected)

	m.CloseConsole(scenario)

	// Close again does not get error.
	m.CloseConsole(scenario)
}
