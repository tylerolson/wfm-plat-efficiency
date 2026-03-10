package standingcalc

import (
	"fmt"
	"sync"
	"time"
)

type Info struct {
	ItemSlug string
	ItemName string
	Err      error
}

type MarketService struct {
	api *marketAPI
}

func NewMarketService(api *marketAPI) *MarketService {
	return &MarketService{
		api: api,
	}
}

// getVendorStats takes in a Vendor which contains a list of the items name and type (mod, weapon, etc).
// It will then call another function to fetch the api, and update the market data.
func (s *MarketService) updateVendorStats(vendor *Vendor) chan Info {
	infoCh := make(chan Info, len(vendor.Items))

	go func() {
		var wg sync.WaitGroup

		ticker := time.NewTicker(time.Second / 3)
		defer ticker.Stop()

		for _, item := range vendor.Items {

			wg.Add(1)

			go func(item *Item) {
				defer wg.Done()

				<-ticker.C

				err := s.UpdateItemStats(item)

				infoCh <- Info{
					item.Slug,
					item.Name,
					err,
				}
			}(item)
		}

		wg.Wait()
		vendor.ScoreItems()
		close(infoCh)
	}()

	return infoCh
}

// UpdateItemStats updates market data for a single item
func (s *MarketService) UpdateItemStats(item *Item) error {
	marketData, err := s.api.GetItemStatistics(item.Slug, item.Type)
	if err != nil {
		return fmt.Errorf("error fetching %s: %w", item.Slug, err)
	}

	item.MarketData = *marketData
	return nil
}
