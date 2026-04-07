package standingcalc

import (
	"encoding/json"
	"fmt"
)

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

func (t ItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *ItemType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try numeric fallback for backward compatibility with JSON files
		var n int
		if err := json.Unmarshal(data, &n); err != nil {
			return err
		}
		*t = ItemType(n)
		return nil
	}
	switch s {
	case "Mod":
		*t = ItemTypeMod
	case "ArchPart":
		*t = ItemTypeArchPart
	case "Weapon":
		*t = ItemTypeWeapon
	default:
		return fmt.Errorf("unknown item type: %s", s)
	}
	return nil
}

type MarketData struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	Score  float64 `json:"score"`
}

type Item struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	Type         ItemType `json:"type"`
	StandingCost int      `json:"standing"`
	MarketData
}

func (i *Item) StandingPerPlat() float64 {
	if i.Price == 0 {
		return 0
	}
	return float64(i.StandingCost) / i.Price
}
