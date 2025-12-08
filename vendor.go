package standingcalc

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"text/tabwriter"
)

type Vendor struct {
	Slug  string  `json:"slug"`
	Name  string  `json:"name"`
	Items []*Item `json:"items"`
}

// MostVolume returns the [Item] with the greatest volume over 2 days
func (v Vendor) MostVolume() *Item {
	if v.Items == nil || len(v.Items) == 0 {
		return nil
	}

	mostVolume := v.Items[0]
	for _, i := range v.Items {
		if i.Volume > mostVolume.Volume {
			mostVolume = i
		}
	}

	return mostVolume
}

// MostProfit returns the [Item] with the greatest weighted platinum average across 2 days
func (v Vendor) MostProfit() *Item {
	if v.Items == nil || len(v.Items) == 0 {
		return nil
	}

	mostProfit := v.Items[0]
	for _, i := range v.Items {
		if i.Price > mostProfit.Price {
			mostProfit = i
		}
	}

	return mostProfit
}

// MostEfficient returns the [Item] with the highest computed score:
// price * ln(1+volume)/ln(1+maxVolume) / standingCost
// This should log scaling to punish low volume greatly but not over score higher volumes.
// We then also normalize based on the vendor's current market.
// It should only return [nil] on an error.
func (v Vendor) MostEfficient() *Item {
	if v.Items == nil || len(v.Items) == 0 {
		return nil
	}

	maxVolumeAmount := v.MostVolume().Volume
	if maxVolumeAmount <= 0 {
		return nil
	}

	denominator := math.Log1p(maxVolumeAmount)

	mostEfficientItem := v.Items[0]
	itemScore := -math.MaxFloat64

	for _, i := range v.Items {
		// prevent / 0 panic
		if i.StandingCost <= 0 {
			continue
		}
		volumeFactor := math.Log1p(i.Volume) / denominator
		score := i.Price * volumeFactor / float64(i.StandingCost)

		if score > itemScore {
			itemScore = score
			mostEfficientItem = i
		} else if score == itemScore {
			// break tie with weighted average
			if i.Price > mostEfficientItem.Price {
				mostEfficientItem = i
			}
		}
	}

	return mostEfficientItem
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
		if l := len(fmt.Sprintf("%.2f", item.Price)); l > maxPrice {
			maxPrice = l
		}
		if l := len(fmt.Sprintf("%.2f", item.Volume)); l > maxVol {
			maxVol = l
		}
		if l := len(fmt.Sprintf("%.2f", float64(item.StandingCost)/item.Price)); l > maxStandVol {
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
			item.Price,
			item.Volume,
			item.StandingPerPlat(),
		)

		if i != len(itemSlice)-1 {
			fmt.Fprintln(w)
		}
	}
	w.Flush()

	return b.String()
}
