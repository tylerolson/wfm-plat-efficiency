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

type Program struct {
	items map[string]struct {
		name            string
		weightedAverage string
	}
}

func newProgram() *Program {
	return &Program{
		items: make(map[string]struct {
			name            string
			weightedAverage string
		}),
	}
}

func (p *Program) getVendorStats(items []Item) error {
	for _, v := range items {
		err := getStatisitics(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func getStatisitics(item Item) error {
	slog.Info("requesting", "name", item.Name, "type", item.Type)

	resp, err := http.Get(fmt.Sprintf("%v/items/%v/statistics", API, item.Name))
	if err != nil {
		return fmt.Errorf("failed to get statistics:%w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	// we could use a struct here, but not every item type has the same fields, and I dont wanna make a bunch of different structs
	var response map[string]any
	err = decoder.Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode statistics:%w", err)
	}

	payload, ok := response["payload"].(map[string]any)
	if !ok {
		slog.Error("could not find payload")
		return errors.New("could not find payload")
	}

	statisticsClosed, ok := payload["statistics_closed"].(map[string]any)
	if !ok {
		slog.Error("could not find statistics_closed")
		return errors.New("could not find statistics_closed")
	}

	statistics90, ok := statisticsClosed["90days"].([]any)
	if !ok {
		slog.Error("could not find 90days")
		return errors.New("could not find 90days")
	}

	// filter rank 0 mods when item is mod
	if item.Type == Mod {
		var mod0 []any

		for _, v := range statistics90 {
			statMap, ok := v.(map[string]any)
			if !ok {
				slog.Error("could not assert statistic entry as map", "error", err)
				return fmt.Errorf("could not assert statistic entry as map: %w", err)
			}

			modRank, ok := statMap["mod_rank"].(float64)
			if !ok {
				slog.Error("could not find mod_rank", "error", err)
				return fmt.Errorf("could not assert statistic entry as map: %w", err)
			}

			if modRank == 0 {
				mod0 = append(mod0, v)
			}
		}

		statistics90 = mod0
	}

	slog.Info("90days", "90days", statistics90[:2])

	return nil
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
