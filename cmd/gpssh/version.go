package main

import (
	"fmt"
	"runtime/debug"
)

var Version string

func init() {
	if Version != "" {
		return
	}
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, s := range buildInfo.Settings {
			// use vcs revision if its available
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				Version = fmt.Sprintf("%s (commit: %s)", buildInfo.Main.Version, s.Value[:7])
				return
			}
		}
		Version = buildInfo.Main.Version
		return
	}
	Version = "(unknown version)"
}
