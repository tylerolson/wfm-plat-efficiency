package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const API = "https://api.warframe.market/v1"

type Program struct {
	vendors map[string]*Vendor
}

func NewProgram() *Program {
	return &Program{
		vendors: make(map[string]*Vendor),
	}
}

func (p *Program) addVendor(vendor Vendor) {
	p.vendors[vendor.Name] = &vendor
}

func LoadVendors(dir string) ([]Vendor, error) {
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

	p := NewProgram()

	vendors, err := LoadVendors("vendors")
	if err != nil {
		slog.Error("failed to get vendors", "error", err)
	}

	for _, v := range vendors {
		p.addVendor(v)
	}

	for _, v := range p.vendors {
		fmt.Printf("Fetching items for: %s\n", v.Name)
		err := v.getVendorStats()
		if err != nil {
			slog.Error("failed to get vendor statistics", "error", err)
		}
	}

	for _, v := range p.vendors {
		fmt.Printf("\n%s:\n", v.Name)
		fmt.Println(v.String())
	}
}
