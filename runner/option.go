package runner

import (
	"fmt"
	"io"
)

type Option func(*Runner)

var sshWords = map[bool]string{
	true:  "yes",
	false: "no",
}

func sshOption(name string, v any) string {
	switch v.(type) {
	case bool:
		v = sshWords[v.(bool)]
	}
	return fmt.Sprintf("-o %v=%v", name, v)
}

// WithColor turns on or off printing any terminal colors
func WithColor(color bool) Option {
	return func(r *Runner) {
		noColor = !color
	}
}

// WithOkColor uses a specific color for printing success
func WithOkColor(color string) Option {
	return func(r *Runner) {
		r.okColor = color
	}
}

// WithFailedColor uses a specific color for printing failures
func WithFailedColor(color string) Option {
	return func(r *Runner) {
		r.failedColor = color
	}
}

// WithVerbose turns off ssh quiet mode
func WithVerbose(verbose bool) Option {
	return func(r *Runner) {
		r.verbose = verbose
	}
}

// WithConnectionAttempts sets the number of connection attempts ssh tries to make
func WithConnectionAttempts(n int) Option {
	return func(r *Runner) {
		r.addOption(sshOption("ConnectionAttempts", n))
	}
}

// WithStrictHostKeyChecking turns on or off ssh StrictHostKeyChecking
func WithStrictHostKeyChecking(b bool) Option {
	return func(r *Runner) {
		r.addOption(sshOption("StrictHostKeyChecking", b))
	}
}

// WithOuput sets the Runner's output writer for both stdout and stderr
func WithOutput(w io.Writer) Option {
	return func(r *Runner) {
		r.logOut = newLogger(w, "")
		r.logErr = newLogger(w, "")
	}
}

// WithInput sets the Runner's input reader
func WithInput(r io.Reader) Option {
	return func(rn *Runner) {
		rn.inputr = r
	}
}
