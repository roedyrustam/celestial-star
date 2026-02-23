package brewer

import (
	"testing"
)

func TestCalculateScore_Espresso(t *testing.T) {
	recipe := Recipe{
		Method:      "espresso",
		CoffeeDose:  18,
		WaterYield:  36,
		Temperature: 91.5,
		BrewTime:    30,
		GrindSize:   "fine",
		RoastLevel:  "medium",
	}

	result := CalculateScore(recipe)

	if result.TotalScore < 95 {
		t.Errorf("Expected high score for ideal espresso, got %d. Feedback: %v", result.TotalScore, result.Feedback)
	}
}

func TestCalculateScore_V60(t *testing.T) {
	recipe := Recipe{
		Method:      "v60",
		CoffeeDose:  15,
		WaterYield:  255, // 1:17
		Temperature: 93,
		BrewTime:    180,
		GrindSize:   "medium-fine",
		RoastLevel:  "light",
	}

	result := CalculateScore(recipe)

	if result.TotalScore < 95 {
		t.Errorf("Expected high score for ideal v60, got %d. Feedback: %v", result.TotalScore, result.Feedback)
	}
}

func TestCalculateScore_FrenchPress(t *testing.T) {
	recipe := Recipe{
		Method:      "frenchpress",
		CoffeeDose:  20,
		WaterYield:  300, // 1:15
		Temperature: 94,
		BrewTime:    270,
		GrindSize:   "coarse",
		RoastLevel:  "medium",
	}

	result := CalculateScore(recipe)

	if result.TotalScore < 95 {
		t.Errorf("Expected high score for ideal french press, got %d. Feedback: %v", result.TotalScore, result.Feedback)
	}
}

func TestCalculateScore_ColdBrew(t *testing.T) {
	recipe := Recipe{
		Method:      "coldbrew",
		CoffeeDose:  50,
		WaterYield:  500, // 1:10
		Temperature: 15,
		BrewTime:    57600, // 16h
		GrindSize:   "coarse",
		RoastLevel:  "dark",
	}

	result := CalculateScore(recipe)

	if result.TotalScore < 95 {
		t.Errorf("Expected high score for ideal cold brew, got %d. Feedback: %v", result.TotalScore, result.Feedback)
	}
}

func TestCalculateScore_BadDose(t *testing.T) {
	recipe := Recipe{
		CoffeeDose: 0,
	}

	result := CalculateScore(recipe)

	if len(result.Feedback) == 0 || result.Feedback[0] != "Dose kopi harus lebih dari 0g." {
		t.Errorf("Expected error feedback for 0 dose, got %v", result.Feedback)
	}
}

func TestPredictFlavor(t *testing.T) {
	// Bitter brew (High temp, long time)
	recipe := Recipe{
		Temperature: 98,
		BrewTime:    300,
		Method:      "v60",
		CoffeeDose:  15,
		WaterYield:  250,
	}
	flavor := predictFlavor(&recipe, 250/15)
	if flavor.Bitterness < 7 {
		t.Errorf("Expected high bitterness for hot/long brew, got %d", flavor.Bitterness)
	}

	// Acidic brew (Low temp, short time)
	recipe2 := Recipe{
		Temperature: 85,
		BrewTime:    120,
		Method:      "v60",
		CoffeeDose:  15,
		WaterYield:  250,
	}
	flavor2 := predictFlavor(&recipe2, 250/15)
	if flavor2.Acidity < 7 {
		t.Errorf("Expected high acidity for cold/short brew, got %d", flavor2.Acidity)
	}
}
