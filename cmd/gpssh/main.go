// Homemade parallel ssh, because why not
package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/kmullin/go-pssh/ssh"
)

func main() {
	var (
		noColor     bool
		showHelp    bool
		showVersion bool
	)

	var (
		fanOut      = 50
		okColor     = "#A8CC8C"
		failedColor = "#E88388"
	)

	// ssh flags
	var (
		connectionAttempts int
		strictHostChecking bool
		verbose            bool
	)

	flags := flag.NewFlagSet("default", flag.ExitOnError)
	flags.BoolVarP(&showHelp, "help", "h", false, "Show help (this output)")
	flags.BoolVarP(&showVersion, "version", "V", false, "Show current version")
	flags.IntVarP(&fanOut, "fanout", "f", fanOut, "Hosts to run in parallel")
	flags.BoolVarP(&noColor, "no-color", "n", false, "Disable colors")
	flags.StringVar(&okColor, "ok-color", okColor, "Color to use for stdout")
	flags.StringVar(&failedColor, "fail-color", failedColor, "Color to use for stderr")
	flags.SortFlags = false

	sshFlags := flag.NewFlagSet("ssh", flag.ExitOnError)
	sshFlags.BoolVarP(&verbose, "verbose", "v", false, "Verbose output (turns off quiet)")
	sshFlags.BoolVarP(&strictHostChecking, "strict", "s", false, "Strict host key checking")
	sshFlags.IntVarP(&connectionAttempts, "retries", "r", 1, "Number of ssh connection attempts")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\t%v [option] command [argument ...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flags.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSSH Options:\n")
		sshFlags.PrintDefaults()
	}

	flag.CommandLine.AddFlagSet(flags)
	flag.CommandLine.AddFlagSet(sshFlags)
	flag.SetInterspersed(false) // this disables interspersed parsing so we dont mess up command to send to SSH
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Need a command to run")
		flag.Usage()
		os.Exit(1)
	}

	runner := ssh.NewRunner(flag.Args(), fanOut,
		ssh.WithColor(!noColor), // disables or enables
		ssh.WithOkColor(okColor),
		ssh.WithFailedColor(failedColor),
		ssh.WithVerbose(verbose),
		ssh.WithConnectionAttempts(connectionAttempts),
		ssh.WithStrictHostKeyChecking(strictHostChecking),

		// ssh.WithOutput(w),
		// ssh.WithInput(r),
	)
	runner.Run()
	runner.SummaryReport()

	if runner.HasErrors() {
		os.Exit(1)
	}
}
