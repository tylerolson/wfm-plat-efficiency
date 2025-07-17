package standingcalc

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
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

type ItemType int

const (
	ItemTypeMod ItemType = iota
	ItemTypeArchPart
	ItemTypeWeapon
)

func (t ItemType) String() string {
	switch t {
	case ItemTypeMod:
		return "Mod"
	case ItemTypeArchPart:
		return "ArchPart"
	case ItemTypeWeapon:
		return "Weapon"
	default:
		return "Unknown"
	}
}

type MarketData struct {
	WeightedAvgPrice float64
	AvgVol           float64
}

type Item struct {
	Name         string   `json:"name"`
	Type         ItemType `json:"type"`
	StandingCost int      `json:"standing"`
	MarketData
}

type Vendor struct {
	Name  string  `json:"name"`
	Items []*Item `json:"items"`
}

func (v Vendor) String() string {
	// get slice from map so we can sort it
	itemSlice := make([]*Item, 0, len(v.Items))

	maxName, maxType, maxStanding, maxPrice, maxVol, maxStandVol := 0, 0, 0, 0, 0, 0
	for _, item := range v.Items {
		itemSlice = append(itemSlice, item)

		if l := len(item.Name); l > maxName {
			maxName = l
		}
		if l := len(item.Type.String()); l > maxType {
			maxType = l
		}
		if l := len(fmt.Sprintf("%d", item.StandingCost)); l > maxStanding {
			maxStanding = l
		}
		if l := len(fmt.Sprintf("%.2f", item.WeightedAvgPrice)); l > maxPrice {
			maxPrice = l
		}
		if l := len(fmt.Sprintf("%.2f", item.AvgVol)); l > maxVol {
			maxVol = l
		}
		if l := len(fmt.Sprintf("%.2f", float64(item.StandingCost)/item.WeightedAvgPrice)); l > maxStandVol {
			maxStandVol = l
		}
	}

	sort.Slice(itemSlice, func(i, j int) bool {
		if itemSlice[i].Type == itemSlice[j].Type {
			return itemSlice[i].Name < itemSlice[j].Name // Sort by Name if Type is the same
		}
		return itemSlice[i].Type < itemSlice[j].Type // Sort by Type first
	})

	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "Name\tType\tStanding\tPrice\tVolume\tStanding/Plat (lower is better)")
	fmt.Fprintf(
		w, "%s\t%s\t%s\t%s\t%s\t%s\n",
		strings.Repeat("-", maxName),
		strings.Repeat("-", maxType),
		strings.Repeat("-", maxStanding),
		strings.Repeat("-", maxPrice),
		strings.Repeat("-", maxVol),
		strings.Repeat("-", maxStandVol),
	)

	for i, item := range itemSlice {
		fmt.Fprintf(
			w, "%s\t%v\t%v\t%.2f\t%.2f\t%0.2f",
			item.Name,
			item.Type.String(),
			item.StandingCost,
			item.WeightedAvgPrice,
			item.AvgVol,
			float64(item.StandingCost)/item.WeightedAvgPrice,
		)

		if i != len(itemSlice)-1 {
			fmt.Fprintln(w)
		}
	}
	w.Flush()

	return b.String()
}

// getVendorStats takes in a Vendor which contains a list of the items name and type (mod, weapon, etc).
// It will then call another function to fetch the api, and update the market data.
func (v *Vendor) GetVendorStats() error {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second / 3)
	errCh := make(chan error, len(v.Items))
	doneCh := make(chan struct{})

	defer ticker.Stop()

	// TODO: add an info channel to get live updates
	for _, item := range v.Items {
		wg.Add(1)

		go func(i *Item) {
			defer wg.Done()

			<-ticker.C

			if err := i.getStatisitics(); err != nil {
				errCh <- fmt.Errorf("error fetching %s: %w", i.Name, err)
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

type ninetyDay struct {
	Volume   int     `json:"volume"`
	AvgPrice float64 `json:"avg_price"`
	ModRank  int     `json:"mod_rank"`
}

type statisticResponse struct {
	Payload struct {
		StatisticsClosed struct {
			NinetyDays []ninetyDay `json:"90days"`
		} `json:"statistics_closed"`
	} `json:"payload"`
}

type apiErrorResponse struct {
	Error string `json:"error"`
}

// getStatisitics takes in an [Item] containing an items name and [ItemType].
// It fetches the API statistics and calculates and returns:
//
//  1. weightedAveragePrice = ((todayAvgPrice * todayVolume) + (yesterdayAvgPrice * yesterdayVolume)) / (todayVolume + yesterdayVolume)
//  2. avgVolume = (todayVolume + yesterdayVolume) / 2
//  3. error
func (i *Item) getStatisitics() error {
	url := fmt.Sprintf("%v/items/%v/statistics", API, i.Name)

	// make a stupid request with an accept header since warframe market redirects without it
	req, err := http.NewRequest("GET", url, nil) // No request body for GET
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get statistics:%w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var response statisticResponse
	err = decoder.Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode statistics:%w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("item does not exist")
		}

		return fmt.Errorf("unkown HTTP error: %v", resp.Status)
	}

	ninetyDays := response.Payload.StatisticsClosed.NinetyDays

	// filter rank 0 mods when item is mod

	if i.Type == ItemTypeMod {
		var mod0 []ninetyDay
		for _, v := range ninetyDays {
			if v.ModRank == 0 {
				mod0 = append(mod0, v)
			}
		}

		response.Payload.StatisticsClosed.NinetyDays = mod0
	}

	today := ninetyDays[0]
	yesterday := ninetyDays[1]

	i.WeightedAvgPrice = 0.0
	i.WeightedAvgPrice += today.AvgPrice*float64(today.Volume) + yesterday.AvgPrice*float64(yesterday.Volume)
	i.WeightedAvgPrice /= float64(today.Volume + yesterday.Volume)

	i.AvgVol = float64(today.Volume+yesterday.Volume) / 2

	return nil
}
