package main

// From https://42bytes.notion.site/WFM-Api-v2-Documentation-5d987e4aa2f74b55a80db1a09932459d

type ItemsResponse struct {
	ApiVersion string      `json:"apiVersion"`
	Data       []ItemShort `json:"data"`
	Error      string      `json:"error"`
}

// ItemShort represents a summary of an item in the game.
type ItemShort struct {
	Id      string   `json:"id"`      // Unique identifier of the item.
	Slug    string   `json:"slug"`    // URL-friendly name of the item
	GameRef string   `json:"gameRef"` // Reference to the item in the game's database.
	Tags    []string `json:"tags"`

	I18n map[string]ItemShortI18n `json:"i18n"` // Localized text for the item in various languages.

	MaxRank        int8    `json:"maxRank,omitempty"`        // Optional maximum rank the item can achieve.
	MaxCharges     int8    `json:"maxCharges,omitempty"`     // Optional maximum chanrges the item can achieve, used for requiem mods.
	Vaulted        *bool   `json:"vaulted,omitempty"`        // Optional flag indicating if the item is vaulted.
	BulkTradable   bool    `json:"bulkTradable,omitempty"`   // Optional flag indicating if the item is bulk tradable.
	Ducats         int16   `json:"ducats,omitempty"`         // Optional Ducats value of the item.
	MaxAmberStars  int8    `json:"maxAmberStars,omitempty"`  // Optional number of amber stars associated with the item.
	MaxCyanStars   int8    `json:"maxCyanStars,omitempty"`   // Optional number of cyan stars associated with the item.
	BaseEndo       int16   `json:"baseEndo,omitempty"`       // Optional base endo value of the item.
	EndoMultiplier float32 `json:"EndoMultiplier,omitempty"` // Optional multiplier for the endo value.
	Subtypes       string  `json:"EndoMultiplier,omitempty"`
}

type ItemShortI18n struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Thumb   string `json:"thumb"`
	SubIcon string `json:"subIcon,omitempty"`
}

type Item struct {
	Id string `json:"id"`

	Slug     string   `json:"slug"`
	GameRef  string   `json:"gameRef"`
	Tags     []string `json:"tags"`
	Tradable bool     `json:"tradable"`

	SetRoot       bool     `json:"setRoot,omitempty"`
	SetParts      []string `json:"setParts,omitempty"`
	QuantityInSet int8     `json:"quantityInSet,omitempty"`

	Rarity       string   `json:"rarity,omitempty"`
	MaxRank      int8     `json:"maxRank,omitempty"`
	MaxCharges   int8     `json:"maxCharges,omitempty"`
	BulkTradable bool     `json:"bulkTradable,omitempty"`
	Subtypes     []string `json:"subtypes,omitempty"`

	MaxAmberStars  int8    `json:"maxAmberStars,omitempty"`
	MaxCyanStars   int8    `json:"maxCyanStars,omitempty"`
	BaseEndo       int16   `json:"baseEndo,omitempty"`
	EndoMultiplier float32 `json:"endoMultiplier,omitempty"`

	Ducats         int16 `json:"ducats,omitempty"`
	ReqMasteryRank *int8 `json:"reqMasteryRank,omitempty"`
	Vaulted        *bool `json:"vaulted,omitempty"`
	TradingTax     int32 `json:"tradingTax,omitempty"`

	I18n map[string]ItemI18n `json:"i18n"`
}

type ItemI18n struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	WikiLink    string `json:"wikiLink,omitempty"`
	Icon        string `json:"icon"`
	Thumb       string `json:"thumb"`
	SubIcon     string `json:"subIcon,omitempty"`
}
