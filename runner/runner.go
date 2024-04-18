package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/muesli/termenv"
)

var noColor = false // default to color

// Runner manages parallel goroutines of ssh workers.
type Runner struct {
	parallel       int         // how many workers to manage in parallel
	hostc          chan string // the input hostname channel
	errc           chan error  // the output error channel
	logOut, logErr *log.Logger // our local loggers used

	inputr io.Reader // our input, defaults to Stdin

	verbose bool // turn off quiet mode

	sshOpts []string // options for every ssh cmd
	sshCmd  []string // command with args to give to ssh

	okColor, failedColor string // colors to use if we have colors enabled
	ok, failed           int    // final tally of ok/failed cmds
}

// New returns an initialized Runner ready to run ssh commands up to parallel at a time.
func New(command []string, parallel int, opts ...Option) *Runner {
	r := &Runner{
		hostc:    make(chan string),
		errc:     make(chan error),
		parallel: parallel,
		sshCmd:   command,
		logOut:   newLogger(os.Stdout, ""),
		logErr:   newLogger(os.Stderr, ""),
		inputr:   os.Stdin,
	}

	for _, opt := range opts {
		opt(r)
	}

	if !r.verbose {
		r.addOption("-q")
	}
	r.addOption(sshOption("PasswordAuthentication", false))
	r.addOption(sshOption("BatchMode", true))
	return r
}

// Run immediately executes the ssh workers and waits for them to complete
func (r *Runner) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(r.parallel)
	for i := 0; i < r.parallel; i++ {
		go func() {
			defer wg.Done()
			for hostname := range r.hostc {
				r.errc <- r.newCmd(ctx, hostname).Run()
			}
		}()
	}

	go func() {
		defer close(r.errc)
		r.readHosts()
		wg.Wait()
	}()

	r.wait() // we wait for our errc to close
}

// wait drains the errc and counting ok/failed from errors returned
func (r *Runner) wait() {
	for err := range r.errc {
		if err != nil {
			r.failed++
			continue
		}
		r.ok++
	}
}

func (r *Runner) HasErrors() bool {
	return r.failed > 0
}

// SummaryReport reads error results from errc to print out a summary of success vs failures
func (r *Runner) SummaryReport() {
	r.logErr.Printf("\ntotal hosts: %v (%v/%v)",
		r.ok+r.failed,
		colorize(r.ok, r.okColor),
		colorize(r.failed, r.failedColor),
	)
}

func (r *Runner) addOption(opt string) {
	r.sshOpts = append(r.sshOpts, opt)
}

// newPrefixLogger sets up a logger for a specific host
func newPrefixLogger(w io.Writer, hostname string, color string) *log.Logger {
	return newLogger(w, fmt.Sprintf("%v: ", colorize(hostname, color)))
}

func newLogger(w io.Writer, prefix string) *log.Logger {
	return log.New(w, prefix, log.Lmsgprefix)
}

// readHosts reads from os.Stdin for hostnames, and sends them to given channel c
// returns a count of the number of hosts processed
func (r *Runner) readHosts() {
	defer close(r.hostc)
	scanner := bufio.NewScanner(r.inputr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		r.hostc <- scanner.Text()
	}
}

func colorize(v any, color string) string {
	if noColor || color == "" {
		return fmt.Sprint(v)
	}
	p := termenv.ColorProfile()
	return termenv.String(fmt.Sprint(v)).Foreground(p.Color(color)).String()
}
