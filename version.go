package main

import (
	"flag"
	"fmt"
)

var (
	Build, Version string
)

var showVersion = flag.Bool("version", false, "Show current version")

// display current major, minor and patch versions (semver)
func version() string {
	return fmt.Sprintf("%s (build: %s)", Version, Build)
}
