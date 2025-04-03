package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

const API = "https://api.warframe.market/v1"

type ItemStat struct {
	Item
	WeightedAvgPrice float64
	AvgVol           float64
}

type Program struct {
	itemStats map[string]ItemStat
}

func NewProgram() *Program {
	return &Program{
		itemStats: make(map[string]ItemStat),
	}
}

func FormatItemStats(items map[string]ItemStat) string {
	// get slice from map so we can sort it
	itemSlice := make([]ItemStat, 0, len(items))

	maxName, maxPrice, maxVol, maxType, maxStanding := 0, 0, 0, 0, 0
	for _, item := range items {
		itemSlice = append(itemSlice, item)
		if len(item.Name) > maxName {
			maxName = len(item.Name)
		}

		priceStr := fmt.Sprintf("%.2f", item.WeightedAvgPrice)
		if len(priceStr) > maxPrice {
			maxPrice = len(priceStr)
		}

		volStr := fmt.Sprintf("%.2f", item.AvgVol)
		if len(volStr) > maxVol {
			maxVol = len(volStr)
		}

		if len(item.Type) > maxType {
			maxType = len(item.Type)
		}

		if len(strconv.Itoa(item.StandingCost)) > maxStanding {
			maxStanding = len(strconv.Itoa(item.StandingCost))
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
	fmt.Fprintln(w, "Name\tPrice\tVolume\tType\tStanding")

	fmt.Fprintf(
		w, "%s\t%s\t%s\t%s\t%s\n",
		strings.Repeat("-", maxName),
		strings.Repeat("-", maxPrice),
		strings.Repeat("-", maxVol),
		strings.Repeat("-", maxType),
		strings.Repeat("-", maxStanding),
	)

	count := 0
	for _, item := range itemSlice {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%v\t%v", item.Name, item.WeightedAvgPrice, item.AvgVol, item.Type, item.StandingCost)
		count++
		if count != len(items) {
			fmt.Fprintln(w)
		}
	}
	w.Flush()

	return b.String()
}

// getVendorStats takes in a Vendor which contains a list of the items name and type (mod, weapon, etc).
// It will then call another function to fetch the api, and save the result to our program.
func (p *Program) getVendorStats(vendor Vendor) error {
	for _, item := range vendor.Items {
		price, vol, err := getStatisitics(item)
		if err != nil {
			return err
		}

		p.itemStats[item.Name] = ItemStat{
			Item:             item,
			WeightedAvgPrice: price,
			AvgVol:           vol,
		}

		slog.Debug("found item", "name", item.Name, "weightedAvgPrice", price, "avgVol", vol)
	}

	slog.Debug("Finished fetching items")

	return nil
}

// getStatisitics takes in an [Item] containing an items name and [ItemType].
// It fetches the API statistics and calculates and returns:
//
//  1. weightedAveragePrice = ((todayAvgPrice * todayVolume) + (yesterdayAvgPrice * yesterdayVolume)) / (todayVolume + yesterdayVolume)
//  2. avgVolume = (todayVolume + yesterdayVolume) / 2
//  3. error
func getStatisitics(item Item) (float64, float64, error) {
	slog.Debug("requesting", "name", item.Name, "type", item.Type)

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
	if item.Type == ItemTypeMod {
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
	debugSet := flag.Bool("debug", false, "enables debug logging")
	flag.Parse()

	level := slog.LevelInfo
	if *debugSet {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	p := NewProgram()

	smallhexis := AribiterOfHexis
	// smallhexis.Items = smallhexis.Items[:5]

	fmt.Println("Fetching items...")
	err := p.getVendorStats(smallhexis)
	if err != nil {
		slog.Error("failed to get vendor statistics", "error", err)
	}

	fmt.Println("\nItems:")
	fmt.Println(FormatItemStats(p.itemStats))
}
