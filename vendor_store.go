package standingcalc

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
)

//go:embed vendors/*.json
var vendorFS embed.FS

type VendorStore struct {
	fs embed.FS
}

func NewVendorStore() *VendorStore {
	return &VendorStore{
		fs: vendorFS,
	}
}

func (r *VendorStore) LoadAllVendors() (map[string]*Vendor, error) {
	files, err := r.fs.ReadDir("vendors")
	if err != nil {
		return nil, fmt.Errorf("error reading embedded vendor directory: %w", err)
	}

	vendors := make(map[string]*Vendor)

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		vendor, err := r.loadVendorFromFile(file.Name())
		if err != nil {
			return nil, fmt.Errorf("error loading vendor from %s: %w", file.Name(), err)
		}

		vendors[vendor.Name] = vendor
	}

	return vendors, nil
}

func (s *VendorStore) loadVendorFromFile(filename string) (*Vendor, error) {
	path := filepath.Join("vendors", filename)
	data, err := s.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded file %s: %w", filename, err)
	}

	var vendor Vendor
	if err := json.Unmarshal(data, &vendor); err != nil {
		return nil, fmt.Errorf("error unmarshaling file %s: %w", filename, err)
	}

	return &vendor, nil
}
