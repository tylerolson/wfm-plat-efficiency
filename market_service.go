package standingcalc

import (
	"fmt"
	"sync"
	"time"
)

// const BASE_URL = "https://api.warframe.market/v1"

type MarketService struct {
	api *MarketAPI
}

func NewMarketService(api *MarketAPI) *MarketService {
	return &MarketService{
		api: api,
	}
}

// getVendorStats takes in a Vendor which contains a list of the items name and type (mod, weapon, etc).
// It will then call another function to fetch the api, and update the market data.
func (s *MarketService) UpdateVendorStats(vendor *Vendor) error {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second / 3)
	errCh := make(chan error, len(vendor.Items))
	doneCh := make(chan struct{})

	defer ticker.Stop()

	// TODO: add an info channel to get live updates
	for _, item := range vendor.Items {
		wg.Add(1)

		go func(item *Item) {
			defer wg.Done()

			<-ticker.C

			err := s.UpdateItemStats(item)
			if err != nil {
				errCh <- fmt.Errorf("error fetching %s: %w", item.Name, err)
				return
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return err
	}
}

// UpdateItemStats updates market data for a single item
func (s *MarketService) UpdateItemStats(item *Item) error {
	marketData, err := s.api.GetItemStatistics(item.Name, item.Type)
	if err != nil {
		return fmt.Errorf("error fetching %s: %w", item.Name, err)
	}

	item.MarketData = *marketData
	return nil
}
