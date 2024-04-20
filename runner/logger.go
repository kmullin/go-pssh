package runner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sync"
)

// consoleMu is a global mutex to protect loggers output on the console
var consoleMu sync.Mutex

type consoleLogger struct {
	l *log.Logger
}

func (cl *consoleLogger) Printf(format string, v ...any) {
	consoleMu.Lock()
	defer consoleMu.Unlock()
	cl.l.Printf(format, v...)
}

func (cl *consoleLogger) Println(v ...any) {
	consoleMu.Lock()
	defer consoleMu.Unlock()
	cl.l.Println(v...)
}

func (cl *consoleLogger) Writer() io.Writer {
	return cl.l.Writer()
}

// scanPrint scans from r and prints the output to the given logger with prefix
// meant to be run in a goroutine, takes a *sync.WaitGroup
func (cl *consoleLogger) scanPrint(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		cl.Println(scanner.Text())
	}
}

// newHostLogger sets up a logger for a specific host
func newHostLogger(w io.Writer, hostname string, color string) *consoleLogger {
	return newLogger(w, fmt.Sprintf("%v: ", colorize(hostname, color)))
}

func newLogger(w io.Writer, prefix string) *consoleLogger {
	return &consoleLogger{log.New(w, prefix, log.Lmsgprefix)}
}
