package standingcalc

import "fmt"

type Scraper struct {
	vendors     map[string]*Vendor
	service     *MarketService
	vendorStore *VendorStore
}

func NewScraper() *Scraper {
	api := NewMarketAPI()
	service := NewMarketService(api)
	store := NewVendorStore()

	return &Scraper{
		vendors:     make(map[string]*Vendor),
		vendorStore: store,
		service:     service,
	}
}

func (s *Scraper) GetVendors() map[string]*Vendor {
	return s.vendors
}

func (s *Scraper) LoadVendors() error {
	vendors, err := s.vendorStore.LoadAllVendors()
	if err != nil {
		return err
	}

	s.vendors = vendors
	return nil
}

// UpdateVendorStats updates market data for a specific vendor
func (s *Scraper) UpdateVendorStats(vendorName string) error {
	vendor, exists := s.vendors[vendorName]
	if !exists {
		return fmt.Errorf("vendor %s not found", vendorName)
	}

	return s.service.UpdateVendorStats(vendor)
}

// UpdateAllVendorStats updates market data for all vendors
func (s *Scraper) UpdateAllVendorStats() error {
	for _, vendor := range s.vendors {
		if err := s.service.UpdateVendorStats(vendor); err != nil {
			return fmt.Errorf("error updating vendor %s: %w", vendor.Name, err)
		}
	}

	return nil
}
