package main

import (
	"fmt"
	"log"

	standingcalc "github.com/tylerolson/wfm-plat-efficiency"
)

func main() {
	scraper := standingcalc.NewScraper()

	if err := scraper.LoadVendors(); err != nil {
		log.Fatal(err)
	}

	// Update market data for all vendors
	if err := scraper.UpdateAllVendorStats(); err != nil {
		log.Fatal(err)
	}

	// Display vendor information
	for _, v := range scraper.GetVendors() {
		fmt.Printf("\n%s:\n", v.Name)
		fmt.Println(v.String())
	}
}
