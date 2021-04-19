package consoledog

import (
	"sync"
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

type session struct {
	console *expect.Console
	state   *vt10x.State
	output  *Buffer
}

// Manager manages console and its state.
type Manager struct {
	sessions map[string]*session
	current  string

	starters []Starter
	closers  []Closer

	test TestingT

	mu sync.Mutex
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
	ctx.Step(`console output matches:`, m.matchConsoleOutput)
}

func (m *Manager) session() *session {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.sessions[m.current]
}

// NewConsole creates a new console.
func (m *Manager) NewConsole(sc *godog.Scenario) (*expect.Console, *vt10x.State) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess := &session{}

	if s, ok := m.sessions[sc.Id]; ok {
		return s.console, s.state
	}

	m.test.Logf("Console: %s (#%s)\n", sc.Name, sc.Id)

	sess.output = new(Buffer)

	console, state, err := vt10x.NewVT10XConsole(expect.WithStdout(sess.output))
	require.NoError(m.test, err)

	sess.console = console
	sess.state = state

	m.sessions[sc.Id] = sess
	m.current = sc.Id

	for _, fn := range m.starters {
		fn(sc, sess.console)
	}

	return sess.console, sess.state
}

// CloseConsole closes the current console.
func (m *Manager) CloseConsole(sc *godog.Scenario) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sess, ok := m.sessions[sc.Id]
	if !ok {
		return
	}

	for _, fn := range m.closers {
		fn(sc)
	}

	m.test.Logf("Raw output: %q\n", sess.output.String())
	// Dump the terminal's screen.
	m.test.Logf("State: \n%s\n", expect.StripTrailingEmptyLines(sess.state.String()))

	delete(m.sessions, sc.Id)
	m.current = ""
}

// Flush flushes console state.
func (m *Manager) Flush() {
	m.session().console.Expect(expect.EOF, expect.PTSClosed, expect.WithTimeout(10*time.Millisecond)) // nolint: errcheck, gosec
}

func (m *Manager) isConsoleOutput(expected *godog.DocString) error {
	m.Flush()

	t := teeError()
	AssertState(t, m.session().state, expected.Content)

	return t.LastError()
}

func (m *Manager) matchConsoleOutput(expected *godog.DocString) error {
	m.Flush()

	t := teeError()
	AssertStateRegex(t, m.session().state, expected.Content)

	return t.LastError()
}

// WithStarter adds a Starter to Manager.
func (m *Manager) WithStarter(s Starter) *Manager {
	m.starters = append(m.starters, s)

	return m
}

// WithCloser adds a Closer to Manager.
func (m *Manager) WithCloser(c Closer) *Manager {
	m.closers = append(m.closers, c)

	return m
}

// New initiates a new console Manager.
func New(t TestingT, options ...Option) *Manager {
	m := &Manager{
		test:     t,
		sessions: make(map[string]*session),
	}

	for _, o := range options {
		o(m)
	}

	return m
}

// WithStarter adds a Starter to Manager.
func WithStarter(s Starter) Option {
	return func(m *Manager) {
		m.WithStarter(s)
	}
}

// WithCloser adds a Closer to Manager.
func WithCloser(c Closer) Option {
	return func(m *Manager) {
		m.WithCloser(c)
	}
}
