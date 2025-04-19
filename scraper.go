package wfmplatefficiency

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
)

const API = "https://api.warframe.market/v1"

//go:embed vendors/*.json
var vendorFS embed.FS

type Scraper struct {
	vendors map[string]*Vendor
}

func NewScraper() *Scraper {
	return &Scraper{
		vendors: make(map[string]*Vendor),
	}
}

func (p *Scraper) GetVendors() map[string]*Vendor {
	return p.vendors
}

func (p *Scraper) LoadVendors() error {
	files, err := vendorFS.ReadDir("vendors")
	if err != nil {
		return fmt.Errorf("error reading embedded vendor directory: %w", err)
	}

	var vendors []Vendor
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		path := filepath.Join("vendors", file.Name())
		data, err := vendorFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading embedded file %s: %w", file.Name(), err)
		}

		var vendor Vendor
		if err := json.Unmarshal(data, &vendor); err != nil {
			return fmt.Errorf("error unmarshaling file %s: %w", file.Name(), err)
		}

		vendors = append(vendors, vendor)
	}

	for _, v := range vendors {
		p.vendors[v.Name] = &v
	}

	return nil
}
