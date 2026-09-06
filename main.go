package main

import (
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/monitoring-forge/flagrun"
)

var version string

const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

type Opt struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

func (o *Opt) Run(_ []string) {
	u := LinuxMemoryPlugin{}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
}

func main() {
	os.Exit(flagrun.Ship(&Opt{}, flagrun.Version(version)))
}
