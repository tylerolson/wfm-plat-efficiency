package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const API = "https://api.warframe.market/v2"

type Program struct {
	items map[string]ItemShort
}

func newProgram() *Program {
	return &Program{items: make(map[string]ItemShort)}
}

func (p *Program) fetchItems() error {
	resp, err := http.Get(fmt.Sprintf("%v/items", API))
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)

	var response ItemsResponse
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	for _, v := range response.Data {
		p.items[v.Slug] = v
	}

	return nil
}

func (p *Program) getVendorItems(vendorList []string) []ItemShort {
	var items = make([]ItemShort, 0)

	if len(p.items) == 0 {
		return nil
	}

	for _, v := range vendorList {
		items = append(items, p.items[v])
	}

	return items
}

func main() {
	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	p := newProgram()
	p.fetchItems()

	hexis := p.getVendorItems(aribiterOfHexis)

	slog.Info(
		"yo",
		"data", hexis,
	)

}
