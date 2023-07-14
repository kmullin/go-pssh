/*
Homemade parallel ssh, because why nott
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

var (
	strictHostChecking = flag.Bool("strict", false, "strict host key checking")
	fanOut             = flag.Int("fanout", 50, "Hosts to run in parallel")
)

var (
	proxyIncludeString = [...]string{".", "/etc/profile.d/proxy.sh", "2>/dev/null", ";"}
)

var (
	logOut = log.New(os.Stdout, "", 0)
	logErr = log.New(os.Stderr, "", 0)
)

func newCmd(host string) *exec.Cmd {
	// ssh defaults
	s := []string{"-q"}
	if !(*strictHostChecking) {
		s = append(s, "-o StrictHostKeyChecking=no")
	}
	s = append(s, "-o PasswordAuthentication=no")
	s = append(s, host)
	s = append(s, proxyIncludeString[0:]...)
	s = append(s, flag.Arg(0))

	cmd := exec.Command("ssh", s...)
	cmd.Stdin = nil
	return cmd
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(version())
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		fmt.Println("Did not provide enough arguments.")
		os.Exit(1)
	}

	var wg sync.WaitGroup

	c := make(chan string)
	errc := make(chan int, 1)

	// spawn ssh workers, reading from c until its closed
	for i := 0; i < *fanOut; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sshWorker(c, errc)
		}()
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	var count int
	for scanner.Scan() {
		c <- scanner.Text()
		count++
	}

	close(c)
	wg.Wait()

	logErr.Printf("\ntotal hosts: %d", count)
	close(errc)
	errn := <-errc
	if errn != 0 {
		os.Exit(1)
	}
}

func sshWorker(c <-chan string, errc chan<- int) {
	for hostname := range c {
		err := doOne(hostname)
		if err != nil {
			select {
			case errc <- 1:
			default:
			}
		}
	}
}

func doOne(hostname string) error {
	cmd := newCmd(hostname)
	prefix := fmt.Sprintf("%s:", hostname)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		logErr.Println(prefix, err)
		return err
	}

	// buffer stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logOut.Println(prefix, scanner.Text())
		}
	}()
	// buffer stderr
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		logErr.Println(prefix, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		if exitCode := checkExitError(err); exitCode > 0 {
			logErr.Println(prefix, "ssh", err)
			return err
		}
	}

	return nil
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
