// Package standingcalc is a tool for calculating the standing efficiency of a vendor
package standingcalc

type Calculator struct {
	service     *MarketService
	vendorStore *VendorStore
}

func NewCalculator() *Calculator {
	api := newMarketAPI()
	service := NewMarketService(api)
	store := newVendorStore()

	return &Calculator{
		vendorStore: store,
		service:     service,
	}
}

func (s *Calculator) LoadVendors() error {
	return s.vendorStore.loadAllVendors()
}

func (s *Calculator) GetVendors() []*Vendor {
	return s.vendorStore.getVendors()
}

func (s *Calculator) GetVendorSlugs() []string {
	return s.vendorStore.getVendorSlugs()
}

func (s *Calculator) GetVendor(slug string) (*Vendor, error) {
	return s.vendorStore.getVendor(slug)
}

// UpdateVendorStats updates market data for a specific vendor
func (s *Calculator) UpdateVendorStats(slug string) (chan Info, error) {
	vendor, err := s.vendorStore.getVendor(slug)
	if err != nil {
		return nil, err
	}

	return s.service.updateVendorStats(vendor), nil
}
