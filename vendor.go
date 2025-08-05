package standingcalc

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

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
			item.StandingPerPlat(),
		)

		if i != len(itemSlice)-1 {
			fmt.Fprintln(w)
		}
	}
	w.Flush()

	return b.String()
}
