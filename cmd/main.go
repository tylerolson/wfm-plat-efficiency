package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	wfmplatefficiency "github.com/tylerolson/wfm-plat-efficiency"
)

func main() {
	debugSet := flag.Bool("debug", false, "enables debug logging")
	flag.Parse()

	level := slog.LevelInfo
	if *debugSet {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	// weird naming, weird way of loading, it just has to work for now
	p := wfmplatefficiency.NewScraper()

	vendors, err := wfmplatefficiency.LoadVendors()
	if err != nil {
		slog.Error("failed to get vendors", "error", err)
	}

	for _, v := range vendors {
		p.AddVendor(v)
	}

	for _, v := range p.Vendors {
		fmt.Printf("Fetching items for: %s\n", v.Name)
		err := v.GetVendorStats()
		if err != nil {
			slog.Error("failed to get vendor statistics", "error", err)
		}
	}

	for _, v := range p.Vendors {
		fmt.Printf("\n%s:\n", v.Name)
		fmt.Println(v.String())
	}
}
