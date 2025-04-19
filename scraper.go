package wfmplatefficiency

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const API = "https://api.warframe.market/v1"

type Scraper struct {
	Vendors map[string]*Vendor
}

func NewScraper() *Scraper {
	return &Scraper{
		Vendors: make(map[string]*Vendor),
	}
}

func (p *Scraper) AddVendor(vendor Vendor) {
	p.Vendors[vendor.Name] = &vendor
}

func LoadVendors() ([]Vendor, error) {
	dir := "vendors"
	var vendors []Vendor

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading vendor directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", file.Name(), err)
		}

		var vendor Vendor
		if err := json.Unmarshal(data, &vendor); err != nil {
			return nil, fmt.Errorf("error unmarshaling file %s: %w", file.Name(), err)
		}

		vendors = append(vendors, vendor)
		slog.Debug("added vendor", "name", vendor.Name)
	}

	return vendors, nil
}
