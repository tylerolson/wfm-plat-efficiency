package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const API = "https://api.warframe.market/v1"

type ItemStat struct {
	name             string
	weightedAvgPrice float64
	avgVol           float64
}

type Program struct {
	items map[string]ItemStat
}

func newProgram() *Program {
	return &Program{
		items: make(map[string]ItemStat),
	}
}

func (p *Program) getVendorStats(items []Item) error {
	for _, v := range items {
		price, vol, err := getStatisitics(v)
		if err != nil {
			return err
		}

		p.items[v.Name] = ItemStat{
			name:             v.Name,
			weightedAvgPrice: price,
			avgVol:           vol,
		}

		slog.Info("found item", "name", v.Name, "weightedAvgPrice", price, "avgVol", vol)
	}

	return nil
}

// getStatisitics returns (weightedAveragePrice, averageVolume, error)
func getStatisitics(item Item) (float64, float64, error) {
	slog.Info("requesting", "name", item.Name, "type", item.Type)

	resp, err := http.Get(fmt.Sprintf("%v/items/%v/statistics", API, item.Name))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get statistics:%w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	// we could use a struct here, but not every item type has the same fields, and I dont wanna make a bunch of different structs
	var response map[string]any
	err = decoder.Decode(&response)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode statistics:%w", err)
	}

	payload, ok := response["payload"].(map[string]any)
	if !ok {
		slog.Error("could not find payload")
		return 0, 0, errors.New("could not find payload")
	}

	statisticsClosed, ok := payload["statistics_closed"].(map[string]any)
	if !ok {
		slog.Error("could not find statistics_closed")
		return 0, 0, errors.New("could not find statistics_closed")
	}

	statistics90, ok := statisticsClosed["90days"].([]any)
	if !ok {
		slog.Error("could not find 90days")
		return 0, 0, errors.New("could not find 90days")
	}

	// filter rank 0 mods when item is mod
	if item.Type == Mod {
		var mod0 []any

		for _, v := range statistics90 {
			statMap, ok := v.(map[string]any)
			if !ok {
				slog.Error("could not assert statistic entry as map", "error", err)
				return 0, 0, fmt.Errorf("could not assert statistic entry as map: %w", err)
			}

			modRank, ok := statMap["mod_rank"].(float64) // json numbers are float64
			if !ok {
				slog.Error("could not find mod_rank", "error", err)
				return 0, 0, fmt.Errorf("could not assert statistic entry as map: %w", err)
			}

			if modRank == 0 {
				mod0 = append(mod0, v)
			}
		}

		statistics90 = mod0
	}

	today := statistics90[0].(map[string]any)
	todayAvgPrice := today["avg_price"].(float64)
	todayVolume := today["volume"].(float64)

	yesterday := statistics90[1].(map[string]any)
	yesterdayAvgPrice := yesterday["avg_price"].(float64)
	yesterdayVolume := yesterday["volume"].(float64)

	weightedAvgPrice := ((todayAvgPrice * todayVolume) + (yesterdayAvgPrice * yesterdayVolume)) / (todayVolume + yesterdayVolume)
	avgVolume := (todayVolume + yesterdayVolume) / 2

	return weightedAvgPrice, avgVolume, nil
}

func main() {
	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	p := newProgram()
	err := p.getVendorStats(AribiterOfHexis)

	if err != nil {
		slog.Error("error", "error", err)
	}
}
