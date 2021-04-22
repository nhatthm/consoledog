# Terminal Emulator for cucumber/godog

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/consoledog)](https://github.com/nhatthm/consoledog/releases/latest)
[![Build Status](https://github.com/nhatthm/consoledog/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/consoledog/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/consoledog/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/consoledog)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/consoledog)](https://goreportcard.com/report/github.com/nhatthm/consoledog)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/consoledog)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

`consoledog` provides a new [`Console`](https://github.com/netflix/go-expect) for each `cucumber/godog` Scenario.

## Prerequisites

- `Go >= 1.14`

## Install

```bash
go get github.com/nhatthm/consoledog
```

## Usage

Initialize a `consoledog.Manager` with `consoledog.New()` then add it into the `ScenarioInitializer`. If you wish to add listeners to `Manager.NewConsole` and
`Manager.CloseConsole` event, use `consoledog.WithStarter` and `consoledog.WithCloser` option in the constructor.

For example:

```go
package mypackage

import (
    "math/rand"
    "testing"

    expect "github.com/Netflix/go-expect"
    "github.com/cucumber/godog"
    "github.com/nhatthm/consoledog"
)

type writer struct {
    console *expect.Console
}

func (w *writer) registerContext(ctx *godog.ScenarioContext) {
    ctx.Step(`write to console:`, func(s *godog.DocString) error {
        _, err := w.console.Write([]byte(s.Content))

        return err
    })
}

func (w *writer) start(_ *godog.Scenario, console *expect.Console) {
    w.console = console
}

func (w *writer) close(_ *godog.Scenario) {
    w.console = nil
}

func TestIntegration(t *testing.T) {
    t.Parallel()

    w := &writer{}
    m := consoledog.New(t,
        consoledog.WithStarter(w.start),
        consoledog.WithCloser(w.close),
    )

    suite := godog.TestSuite{
        Name: "Integration",
        ScenarioInitializer: func(ctx *godog.ScenarioContext) {
            m.RegisterContext(ctx)
        },
        Options: &godog.Options{
            Strict:    true,
            Output:    out,
            Randomize: rand.Int63(),
        },
    }

    // Run the suite.
}
```

See more: [#Examples](#Examples)

### Resize

In case you want to resize the terminal (default is `80x100`) to avoid text wrapping, for example:

```gherkin
        Then console output is:
        """
        panic: could not build credentials provider option: unsupported credentials prov
        ider
        """
```

Use `consoledog.WithTermSize(cols, rows)` while initiating with `consoledog.New()`, for example:

```go
package mypackage

import (
    "testing"

    "github.com/nhatthm/consoledog"
)

func TestIntegration(t *testing.T) {
    // ...
    m := consoledog.New(t, consoledog.WithTermSize(100, 100))
    // ...
}
```

Then your step will become:

```gherkin
        Then console output is:
        """
        panic: could not build credentials provider option: unsupported credentials provider
        """
```

## Steps

### `console output is:`

Asserts the output of the console.

For example:

```gherkin
    Scenario: Find all transaction in range with invalid format
        When I run command "transactions -d --format invalid"

        Then console output is:
        """
        panic: unknown output format
        """
```

## Examples

Full suite: https://github.com/nhatthm/consoledog/tree/master/features

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
