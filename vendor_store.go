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
	vendors []*Vendor
	fs      embed.FS
}

func newVendorStore() *VendorStore {
	return &VendorStore{
		vendors: make([]*Vendor, 0),
		fs:      vendorFS,
	}
}

func (s *VendorStore) getVendors() []*Vendor {
	return s.vendors
}

func (s *VendorStore) getVendor(name string) (*Vendor, error) {
	for _, v := range s.vendors {
		if v.Name == name {
			return v, nil
		}
	}

	return nil, fmt.Errorf("vendor %s does not exist", name)
}

func (s *VendorStore) loadAllVendors() error {
	files, err := s.fs.ReadDir("vendors")
	if err != nil {
		return fmt.Errorf("error reading embedded vendor directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		vendor, err := s.loadVendorFromFile(file.Name())
		if err != nil {
			return fmt.Errorf("error loading vendor from %s: %w", file.Name(), err)
		}

		s.vendors = append(s.vendors, vendor)

	}

	return nil
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
