package standingcalc

type Scraper struct {
	service     *MarketService
	vendorStore *VendorStore
}

func NewScraper() *Scraper {
	api := newMarketAPI()
	service := NewMarketService(api)
	store := newVendorStore()

	return &Scraper{
		vendorStore: store,
		service:     service,
	}
}

func (s *Scraper) LoadVendors() error {
	return s.vendorStore.loadAllVendors()
}

func (s *Scraper) GetVendors() []*Vendor {
	return s.vendorStore.getVendors()
}

func (s *Scraper) GetVendor(name string) (*Vendor, error) {
	return s.vendorStore.getVendor(name)
}

// UpdateVendorStats updates market data for a specific vendor
func (s *Scraper) UpdateVendorStats(vendorName string) (chan Info, error) {
	vendor, err := s.vendorStore.getVendor(vendorName)
	if err != nil {
		return nil, err
	}

	return s.service.updateVendorStats(vendor), nil
}
