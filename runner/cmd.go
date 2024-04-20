package runner

import (
	"context"
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
		logOut: newPrefixLogger(r.logOut.Writer(), hostname, r.okColor),
		logErr: newPrefixLogger(r.logErr.Writer(), hostname, r.failedColor),
	}

	var args []string
	args = append(args, r.sshOpts...)
	args = append(args, hostname)
	args = append(args, preamble...)
	args = append(args, r.sshCmd...)

	sc.cmd = exec.CommandContext(ctx, "ssh", args...)
	sc.cmd.Stdin = nil
	sc.cmd.Cancel = func() error {
		return sc.cmd.Process.Signal(syscall.SIGINT)
	}

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

	if err = sc.cmd.Start(); err != nil {
		sc.logErr.Println(err)
		return
	}

	// handle our output, wait for all printing to be done
	var wg sync.WaitGroup
	wg.Add(2)
	go sc.logOut.scanPrint(&wg, stdout)
	go sc.logErr.scanPrint(&wg, stderr)
	// incorrect to call cmd.Wait before all reads from the pipe have completed
	// so we wait on all reads to complete first
	wg.Wait()

	err = sc.cmd.Wait()
	if err != nil {
		sc.logErr.Println("ssh:", err)
	}

	return
}
