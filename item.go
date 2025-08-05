package standingcalc

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

func (i *Item) StandingPerPlat() float64 {
	if i.WeightedAvgPrice == 0 {
		return 0
	}
	return float64(i.StandingCost) / i.WeightedAvgPrice
}
