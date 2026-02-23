// brewer/models.go
package brewer

type Recipe struct {
	Method      string  `json:"method"`      // e.g., "v60", "espresso", "aeropress"
	CoffeeDose  float64 `json:"coffee_dose"` // grams
	WaterYield  float64 `json:"water_yield"` // grams
	Temperature float64 `json:"temperature"` // Celsius
	BrewTime    int     `json:"brew_time"`   // seconds
	GrindSize   string  `json:"grind_size"`  // "fine", "medium-fine", "medium", "medium-coarse", "coarse"
	RoastLevel  string  `json:"roast_level"` // "light", "medium", "dark"
}

type FlavorProfile struct {
	Acidity    int `json:"acidity"`    // 1-10
	Sweetness  int `json:"sweetness"`  // 1-10
	Body       int `json:"body"`       // 1-10
	Bitterness int `json:"bitterness"` // 1-10
}

type ScoreResult struct {
	TotalScore      int           `json:"total_score"`
	RatioScore      int           `json:"ratio_score"` // out of 40
	TempScore       int           `json:"temp_score"`  // out of 30
	TimeScore       int           `json:"time_score"`  // out of 30
	Feedback        []string      `json:"feedback"`
	IdealRatio      string        `json:"ideal_ratio"`
	IdealTemp       string        `json:"ideal_temp"`
	IdealTime       string        `json:"ideal_time"`
	CalculatedRatio float64       `json:"calculated_ratio"`
	Flavor          FlavorProfile `json:"flavor"`
}
