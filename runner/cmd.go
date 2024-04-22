package runner

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
	"syscall"
)

// preamble is used to source proxy variables if they exist
// this is used in environemnts where you might need to use an HTTP proxy
var preamble = []string{".", "/etc/profile.d/proxy.sh", "2>/dev/null", ";"}

// cmd executes ssh using loggers for input/output
type cmd struct {
	cmd *exec.Cmd

	logOut, logErr *consoleLogger
}

func (r *Runner) newCmd(ctx context.Context, hostname string) *cmd {
	sc := &cmd{
		logOut: newHostLogger(r.logOut.Writer(), hostname, r.okColor),
		logErr: newHostLogger(r.logErr.Writer(), hostname, r.failedColor),
	}

	var args []string
	args = append(args, r.sshOpts...)
	args = append(args, hostname)
	args = append(args, preamble...)
	args = append(args, r.sshCmd...)

	sc.cmd = exec.CommandContext(ctx, "ssh", args...)
	sc.cmd.Stdin = nil
	if sc.cmd.SysProcAttr == nil {
		sc.cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	sc.cmd.SysProcAttr.Setpgid = true

	return sc
}

// Run wraps the exec.Cmd Run method with StdOut/Err logging.
func (sc *cmd) Run() (err error) {
	stdout, err := sc.cmd.StdoutPipe()
	if err != nil {
		sc.logErr.Println(err)
		return
	}
	stderr, err := sc.cmd.StderrPipe()
	if err != nil {
		sc.logErr.Println(err)
		return
	}

	sc.cmd.Cancel = func() error {
		// force close our WriteCloser so our scanners stop immediately
		defer stdout.Close()
		defer stderr.Close()
		return sc.cmd.Process.Signal(syscall.SIGINT)
	}

	if err = sc.cmd.Start(); err != nil {
		sc.logErr.Println(err)
		return
	}

	// handle our output, wait for all printing to be done
	var wg sync.WaitGroup
	wg.Add(2)
	go scanPrint(&wg, stdout, sc.logOut)
	go scanPrint(&wg, stderr, sc.logErr)
	// incorrect to call cmd.Wait before all reads from the pipe have completed
	// so we wait on all reads to complete first
	wg.Wait()

	if err = sc.cmd.Wait(); err != nil {
		sc.logErr.Println("ssh:", err)
	}

	return
}

// scanPrint scans from r and prints the output to the given logger with prefix
// meant to be run in a goroutine, takes a *sync.WaitGroup
func scanPrint(wg *sync.WaitGroup, r io.Reader, l *consoleLogger) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l.Println(scanner.Text())
	}
}
