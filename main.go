package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/prometheus/procfs"
)

// version by Makefile
var version string

type cmdOpts struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

type LinuxMemoryPlugin struct{}

func (u LinuxMemoryPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"": {
			Label: "Linux Memory",
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "total", Label: "Total", Stacked: false},
				{Name: "available", Label: "Available", Stacked: false},
				{Name: "used", Label: "Used", Stacked: false},
				{Name: "kernelstack", Label: "KernelStack", Stacked: true},
				{Name: "vmallocused", Label: "VmallocUsed", Stacked: true},
				{Name: "pagetables", Label: "PageTables", Stacked: true},
				{Name: "mapped", Label: "Mapped", Stacked: true},
				{Name: "anonpages", Label: "AnonPages", Stacked: true},
				{Name: "slab", Label: "Slab", Stacked: true},
				{Name: "buffers", Label: "Buffers", Stacked: true},
				{Name: "cached", Label: "Cached", Stacked: true},
				{Name: "free", Label: "Free", Stacked: true},
			},
		},
	}
}

func (u LinuxMemoryPlugin) MetricKeyPrefix() string {
	return "linux-memory"
}

func meminfoValue(v *uint64) float64 {
	if v == nil {
		return 0
	}
	return float64(*v * 1024)
}

func meminfoSub(total *uint64, subs ...*uint64) float64 {
	if total == nil {
		return 0
	}
	remaining := *total
	for _, v := range subs {
		if v != nil {
			remaining -= *v
		}
	}
	return float64(remaining * 1024)
}

func (u LinuxMemoryPlugin) FetchMetrics() (map[string]float64, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return nil, err
	}
	m, err := fs.Meminfo()
	if err != nil {
		return nil, err
	}

	result := map[string]float64{
		"total":       meminfoValue(m.MemTotal),
		"kernelstack": meminfoValue(m.KernelStack),
		"vmallocused": meminfoValue(m.VmallocUsed),
		"pagetables":  meminfoValue(m.PageTables),
		"mapped":      meminfoValue(m.Mapped),
		"anonpages":   meminfoValue(m.AnonPages),
		"slab":        meminfoValue(m.Slab),
		"buffers":     meminfoValue(m.Buffers),
		"cached":      meminfoValue(m.Cached),
		"free":        meminfoValue(m.MemFree),
	}

	if m.MemAvailable != nil {
		result["used"] = meminfoSub(m.MemTotal, m.MemAvailable)
		result["available"] = meminfoValue(m.MemAvailable)
	} else {
		result["used"] = meminfoSub(m.MemTotal, m.MemFree, m.Buffers, m.Cached)
	}

	return result, nil
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := cmdOpts{}
	psr := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opts.Version {
		fmt.Printf(`%s %s
Compiler: %s %s
`,
			os.Args[0],
			version,
			runtime.Compiler,
			runtime.Version())
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	u := LinuxMemoryPlugin{}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Run()
	return 0
}
