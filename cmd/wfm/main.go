package main

import (
	"fmt"
	"log/slog"
	"os"

	standingcalc "github.com/tylerolson/wfm-plat-efficiency"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// weird naming, weird way of loading, it just has to work for now
	scraper := standingcalc.NewScraper()

	err := scraper.LoadVendors()
	if err != nil {
		slog.Error("failed to get vendors", "error", err)
		return
	}

	for _, v := range scraper.GetVendors() {
		fmt.Printf("Fetching items for: %s\n", v.Name)
		err := v.GetVendorStats()
		// TODO: get live updates via channel
		if err != nil {
			slog.Error("failed to get vendor statistics", "error", err)
		}
	}

	for _, v := range scraper.GetVendors() {
		fmt.Printf("\n%s:\n", v.Name)
		fmt.Println(v.String())
	}
}
