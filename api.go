package standingcalc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const API = "https://api.warframe.market/v1"

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

// MarketAPI handles all API communication
type MarketAPI struct {
	client  *http.Client
	baseURL string
}

func NewMarketAPI() *MarketAPI {
	return &MarketAPI{
		client:  &http.Client{},
		baseURL: API,
	}
}

// getStatisitics takes in an [Item] containing an items name and [ItemType].
// It fetches the API statistics and calculates and returns:
//
//  1. weightedAveragePrice = ((todayAvgPrice * todayVolume) + (yesterdayAvgPrice * yesterdayVolume)) / (todayVolume + yesterdayVolume)
//  2. avgVolume = (todayVolume + yesterdayVolume) / 2
//  3. error
func (api *MarketAPI) GetItemStatistics(itemName string, itemType ItemType) (*MarketData, error) {
	url := fmt.Sprintf("%v/items/%v/statistics", api.baseURL, itemName)

	// make a stupid request with an accept header since warframe market redirects without it
	req, err := http.NewRequest("GET", url, nil) // No request body for GET
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics:%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("item does not exist")
		}

		return nil, fmt.Errorf("unkown HTTP error: %v", resp.Status)
	}

	var response statisticResponse
	if err != json.NewDecoder(resp.Body).Decode(&response) {
		return nil, fmt.Errorf("failed to decode statistics:%w", err)
	}

	return api.calculateMarketData(response.Payload.StatisticsClosed.NinetyDays, itemType)
}

func (api *MarketAPI) calculateMarketData(ninetyDays []ninetyDay, itemType ItemType) (*MarketData, error) {
	if len(ninetyDays) < 2 {
		return nil, fmt.Errorf("insufficient data points")
	}

	// filter rank 0 mods when item is mod
	if itemType == ItemTypeMod {
		var mod0 []ninetyDay
		for _, v := range ninetyDays {
			if v.ModRank == 0 {
				mod0 = append(mod0, v)
			}
		}

		ninetyDays = mod0

		if len(ninetyDays) < 2 {
			return nil, fmt.Errorf("insufficient rank 0 mod data")
		}
	}

	today := ninetyDays[0]
	yesterday := ninetyDays[1]

	totalVolume := today.Volume + yesterday.Volume
	if totalVolume == 0 {
		return nil, fmt.Errorf("no trading volume data")
	}

	weightedAvgPrice := (today.AvgPrice*float64(today.Volume) + yesterday.AvgPrice*float64(yesterday.Volume)) / float64(totalVolume)
	avgVol := float64(totalVolume) / 2

	return &MarketData{
		WeightedAvgPrice: weightedAvgPrice,
		AvgVol:           avgVol,
	}, nil
}
