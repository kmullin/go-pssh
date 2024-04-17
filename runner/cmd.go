package runner

import (
	"bufio"
	"io"
	"log"
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

	logOut, logErr *log.Logger
}

func (r *Runner) newCmd(hostname string) *cmd {
	sc := &cmd{
		logOut: newPrefixLogger(r.logOut.Writer(), hostname, r.okColor),
		logErr: newPrefixLogger(r.logErr.Writer(), hostname, r.failedColor),
	}

	var args []string
	args = append(args, r.sshOpts...)
	args = append(args, hostname)
	args = append(args, preamble...)
	args = append(args, r.sshCmd...)

	sc.cmd = exec.Command("ssh", args...)
	sc.cmd.Stdin = nil
	return sc
}

func (sc *cmd) Run() error {
	stdout, err := sc.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := sc.cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := sc.cmd.Start(); err != nil {
		sc.logErr.Println(err)
		return err
	}

	// handle our output, wait for all printing to be done
	var wg sync.WaitGroup
	wg.Add(2)
	go scanPrint(&wg, stdout, sc.logOut)
	go scanPrint(&wg, stderr, sc.logErr)
	wg.Wait()

	if err := sc.cmd.Wait(); err != nil {
		if exitCode := checkExitError(err); exitCode > 0 {
			sc.logErr.Println("ssh:", err)
			return err
		}
	}

	return nil
}

// scanPrint scans from r and prints the output to the given logger with prefix
func scanPrint(wg *sync.WaitGroup, r io.Reader, l *log.Logger) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l.Println(scanner.Text())
	}
}

// checkExitError will check if error has an ExitStatus and returns that status code.
func checkExitError(err error) int {
	// did command return an exit code > 0
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}
