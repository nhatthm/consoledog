package consoledog

import (
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/cucumber/godog"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

// Starter is a callback when console starts.
type Starter func(sc *godog.Scenario, console *expect.Console)

// Closer is a callback when console closes.
type Closer func(sc *godog.Scenario)

// Option configures Manager.
type Option func(m *Manager)

// Manager manages console and its state.
type Manager struct {
	test *testing.T

	console *expect.Console
	state   *vt10x.State
	output  *Buffer

	starters []Starter
	closers  []Closer
}

type tHelper interface {
	Helper()
}

// RegisterContext register console Manager to test context.
func (m *Manager) RegisterContext(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		m.NewConsole(sc)
	})

	ctx.AfterScenario(func(sc *godog.Scenario, _ error) {
		m.CloseConsole(sc)
	})

	ctx.Step(`console output is:`, m.isConsoleOutput)
}

// NewConsole creates a new console.
func (m *Manager) NewConsole(sc *godog.Scenario) (*expect.Console, *vt10x.State) {
	if m.console != nil {
		return m.console, m.state
	}

	m.test.Logf("Console: %s (#%s)\n", sc.Name, sc.Id)

	m.output = new(Buffer)

	console, state, err := vt10x.NewVT10XConsole(expect.WithStdout(m.output))
	require.NoError(m.test, err)

	m.console = console
	m.state = state

	for _, fn := range m.starters {
		fn(sc, console)
	}

	return console, state
}

// CloseConsole closes the current console.
func (m *Manager) CloseConsole(sc *godog.Scenario) {
	if m.console == nil {
		return
	}

	for _, fn := range m.closers {
		fn(sc)
	}

	m.test.Logf("Raw output: %q\n", m.output.String())
	// Dump the terminal's screen.
	m.test.Logf("State: \n%s\n", expect.StripTrailingEmptyLines(m.state.String()))

	m.console = nil
	m.state = nil
	m.output = nil
}

// Flush flushes console state.
func (m *Manager) Flush() {
	m.console.Expect(expect.EOF, expect.PTSClosed, expect.WithTimeout(10*time.Millisecond)) // nolint: errcheck, gosec
}

func (m *Manager) isConsoleOutput(expected *godog.DocString) error {
	m.Flush()

	t := t()
	AssertState(t, m.state, expected.Content)

	return t.LastError()
}

// New initiates a new console Manager.
func New(t *testing.T, options ...Option) *Manager { // nolint: thelper
	m := &Manager{
		test: t,
	}

	for _, o := range options {
		o(m)
	}

	return m
}

// WithStarter adds a Starter to Manager.
func WithStarter(s Starter) Option {
	return func(m *Manager) {
		m.starters = append(m.starters, s)
	}
}

// WithCloser adds a Closer to Manager.
func WithCloser(c Closer) Option {
	return func(m *Manager) {
		m.closers = append(m.closers, c)
	}
}
