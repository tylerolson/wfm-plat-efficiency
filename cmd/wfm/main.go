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
	for _, vendor := range scraper.GetVendors() {
		resultChan, err := scraper.UpdateVendorStats(vendor.Name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Starting %v\n", vendor.Name)
		for value := range resultChan { // Loop until the channel is closed
			if value.Err != nil {
				fmt.Printf("Failed to fetch %s: %v\n", value.ItemName, value.Err)
			} else {
				fmt.Printf("Fetched %v\n", value.ItemName)
			}
		}

	}

	// Display vendor information
	for _, v := range scraper.GetVendors() {
		fmt.Printf("\n%s:\n", v.Name)
		fmt.Println(v.String())
	}
}
