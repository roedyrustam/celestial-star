// brewer/engine.go
package brewer

import (
	"fmt"
	"math"
)

// CalculateScore determines the quality of a brew based on specialty coffee standards
func CalculateScore(r Recipe) ScoreResult {
	if r.CoffeeDose <= 0 {
		return ScoreResult{Feedback: []string{"Dose kopi harus lebih dari 0g."}}
	}

	ratio := r.WaterYield / r.CoffeeDose

	result := ScoreResult{
		CalculatedRatio: ratio,
		Feedback:        []string{},
	}

	switch r.Method {
	case "espresso":
		evaluateEspresso(&r, ratio, &result)
	case "v60", "pourover":
		evaluatePourover(&r, ratio, &result)
	case "aeropress":
		evaluateAeropress(&r, ratio, &result)
	case "frenchpress":
		evaluateFrenchPress(&r, ratio, &result)
	case "coldbrew":
		evaluateColdBrew(&r, ratio, &result)
	default:
		// Default fallback to pourover
		evaluatePourover(&r, ratio, &result)
	}

	result.TotalScore = result.RatioScore + result.TempScore + result.TimeScore

	if result.TotalScore >= 95 {
		result.Feedback = append([]string{"Luar biasa! Seduhan yang sangat ideal dan seimbang."}, result.Feedback...)
	} else if result.TotalScore >= 80 {
		result.Feedback = append([]string{"Seduhan yang bagus, ada sedikit ruang untuk perbaikan."}, result.Feedback...)
	} else {
		result.Feedback = append([]string{"Parameter seduhan perlu disesuaikan untuk hasil yang optimal."}, result.Feedback...)
	}

	result.Flavor = predictFlavor(&r, ratio)

	return result
}

func evaluateEspresso(r *Recipe, ratio float64, res *ScoreResult) {
	res.IdealRatio = "1:2 - 1:2.5"
	res.IdealTemp = "90 - 93 °C"
	res.IdealTime = "25 - 35 detik"

	// Ratio (Espresso is generally 1:2 to 1:2.5)
	if ratio >= 1.8 && ratio <= 2.8 {
		res.RatioScore = 40
		if ratio < 2.0 {
			res.Feedback = append(res.Feedback, "Rasio Ristretto (1:<2), akan menghasilkan bodi tebal dan rasa intens.")
		} else if ratio > 2.5 {
			res.Feedback = append(res.Feedback, "Rasio Lungo (1:>2.5), akan menghasilkan ekstraksi lebih tinggi, berisiko over-ekstraksi.")
		}
	} else {
		res.RatioScore = max(0, 40-int(math.Abs(ratio-2.2)*15))
		res.Feedback = append(res.Feedback, fmt.Sprintf("Rasio %.1f kurang ideal untuk Espresso. Usahakan mendekati 1:2 hingga 1:2.5.", ratio))
	}

	// Temperature
	if r.Temperature >= 90 && r.Temperature <= 93 {
		res.TempScore = 30
	} else {
		res.TempScore = max(0, 30-int(math.Abs(r.Temperature-91.5)*5))
		if r.Temperature < 90 {
			res.Feedback = append(res.Feedback, "Suhu terlalu rendah, berisiko under-ekstraksi (asam/sour).")
		} else {
			res.Feedback = append(res.Feedback, "Suhu terlalu tinggi, berisiko over-ekstraksi (pahit/harsh).")
		}
	}

	// Time
	if r.BrewTime >= 25 && r.BrewTime <= 35 {
		res.TimeScore = 30
	} else {
		res.TimeScore = max(0, 30-int(math.Abs(float64(r.BrewTime)-30)*2))
		if r.BrewTime < 25 {
			res.Feedback = append(res.Feedback, "Waktu ekstraksi terlalu cepat. Coba haluskan ukuran gilingan (grind size).")
		} else {
			res.Feedback = append(res.Feedback, "Waktu ekstraksi terlalu lama. Coba kasarkan ukuran gilingan (grind size).")
		}
	}
}

func evaluatePourover(r *Recipe, ratio float64, res *ScoreResult) {
	res.IdealRatio = "1:15 - 1:17"
	res.IdealTemp = "90 - 96 °C"
	res.IdealTime = "150 - 210 detik (2.5 - 3.5 menit)"

	// Ratio (Pourover uses WaterYield as total water poured, usually 1:15 to 1:17)
	if ratio >= 15.0 && ratio <= 17.0 {
		res.RatioScore = 40
	} else {
		res.RatioScore = max(0, 40-int(math.Abs(ratio-16.0)*4))
		if ratio < 15.0 {
			res.Feedback = append(res.Feedback, fmt.Sprintf("Rasio %.1f (1:%.1f) cenderung menghasilkan kopi yang terlalu kuat/pekat.", ratio, ratio))
		} else {
			res.Feedback = append(res.Feedback, fmt.Sprintf("Rasio %.1f (1:%.1f) cenderung menghasilkan kopi yang watery/encer.", ratio, ratio))
		}
	}

	// Temperature (Depends a bit on roast, but generally 90-96)
	targetTemp := 93.0
	if r.RoastLevel == "light" {
		targetTemp = 95.0 // hotter for light
	} else if r.RoastLevel == "dark" {
		targetTemp = 88.0 // cooler for dark
	}

	if math.Abs(r.Temperature-targetTemp) <= 2 {
		res.TempScore = 30
	} else {
		res.TempScore = max(0, 30-int(math.Abs(r.Temperature-targetTemp)*3))
		if r.Temperature < targetTemp {
			res.Feedback = append(res.Feedback, fmt.Sprintf("Suhu terlalu rendah untuk roasting %s. Coba naikkan ke ~%.0f°C.", r.RoastLevel, targetTemp))
		} else {
			res.Feedback = append(res.Feedback, fmt.Sprintf("Suhu terlalu tinggi untuk roasting %s. Coba turunkan ke ~%.0f°C untuk menghindari kepahitan berlebih.", r.RoastLevel, targetTemp))
		}
	}

	// Time
	if r.BrewTime >= 150 && r.BrewTime <= 210 {
		res.TimeScore = 30
	} else {
		res.TimeScore = max(0, 30-int(math.Abs(float64(r.BrewTime)-180)*0.5))
		if r.BrewTime < 150 {
			res.Feedback = append(res.Feedback, "Waktu seduh terlalu cepat. Periksa ukuran gilingan (mungkin terlalu kasar) atau tuangan air (terlalu cepat).")
		} else {
			res.Feedback = append(res.Feedback, "Waktu seduh terlalu lama. Periksa ukuran gilingan (mungkin terlalu halus) yang menyebabkan genangan (stalling).")
		}
	}
}

func evaluateAeropress(r *Recipe, ratio float64, res *ScoreResult) {
	res.IdealRatio = "1:10 - 1:14"
	res.IdealTemp = "80 - 90 °C"
	res.IdealTime = "90 - 150 detik"

	// Ratio
	if ratio >= 10.0 && ratio <= 14.0 {
		res.RatioScore = 40
	} else {
		res.RatioScore = max(0, 40-int(math.Abs(ratio-12.0)*5))
		res.Feedback = append(res.Feedback, fmt.Sprintf("Rasio %.1f kurang lazim untuk Aeropress. Idealnya 1:10 hingga 1:14 (tergantung bypass).", ratio))
	}

	// Temperature
	if r.Temperature >= 80 && r.Temperature <= 90 {
		res.TempScore = 30
	} else {
		res.TempScore = max(0, 30-int(math.Abs(r.Temperature-85.0)*2))
		if r.Temperature > 90 {
			res.Feedback = append(res.Feedback, "Aeropress seringkali optimal di suhu yang lebih rendah (80-90°C) untuk menonjolkan sweetness.")
		}
	}

	// Time
	if r.BrewTime >= 90 && r.BrewTime <= 150 {
		res.TimeScore = 30
	} else {
		res.TimeScore = max(0, 30-int(math.Abs(float64(r.BrewTime)-120)*1))
		res.Feedback = append(res.Feedback, "Waktu steep/plunge di luar kisaran tipikal (90-150s). Perhatikan keseimbangan ekstraksi.")
	}
}

func evaluateFrenchPress(r *Recipe, ratio float64, res *ScoreResult) {
	res.IdealRatio = "1:14 - 1:16"
	res.IdealTemp = "92 - 96 °C"
	res.IdealTime = "240 - 300 detik (4-5 menit)"

	// Ratio
	if ratio >= 14.0 && ratio <= 16.0 {
		res.RatioScore = 40
	} else {
		res.RatioScore = max(0, 40-int(math.Abs(ratio-15.0)*5))
		res.Feedback = append(res.Feedback, fmt.Sprintf("Rasio %.1f. French Press biasanya optimal di 1:15 untuk keseimbangan bodi dan kejernihan.", ratio))
	}

	// Temperature
	if r.Temperature >= 92 && r.Temperature <= 96 {
		res.TempScore = 30
	} else {
		res.TempScore = max(0, 30-int(math.Abs(r.Temperature-94.0)*2))
		if r.Temperature < 92 {
			res.Feedback = append(res.Feedback, "Suhu terlalu rendah. French Press butuh panas untuk ekstraksi perendaman (immersion) yang optimal.")
		}
	}

	// Time
	if r.BrewTime >= 240 && r.BrewTime <= 300 {
		res.TimeScore = 30
	} else {
		res.TimeScore = max(0, 30-int(math.Abs(float64(r.BrewTime)-270)*0.2))
		if r.BrewTime < 240 {
			res.Feedback = append(res.Feedback, "Waktu seduh terlalu singkat untuk metode immersion. Kopi mungkin terasa kurang 'sweet' (underdev).")
		} else {
			res.Feedback = append(res.Feedback, "Waktu seduh sangat lama. Berisiko mengekstraksi rasa pahit yang tidak diinginkan.")
		}
	}

	// Grind Size Check
	if r.GrindSize != "coarse" {
		res.TotalScore -= 5
		res.Feedback = append(res.Feedback, "Peringatan: French Press sebaiknya menggunakan gilingan Coarse (Kasar) untuk menghindari endapan berlebih.")
	}
}

func evaluateColdBrew(r *Recipe, ratio float64, res *ScoreResult) {
	res.IdealRatio = "1:8 - 1:12"
	res.IdealTemp = "4 - 25 °C"
	res.IdealTime = "12 - 24 jam (43200 - 86400 detik)"

	// Ratio (Cold brew is often a concentrate)
	if ratio >= 8.0 && ratio <= 12.0 {
		res.RatioScore = 40
	} else {
		res.RatioScore = max(0, 40-int(math.Abs(ratio-10.0)*3))
		res.Feedback = append(res.Feedback, "Rasio Cold Brew biasanya lebih pekat (1:8 - 1:12) karena ekstraksi dingin lebih lambat.")
	}

	// Temperature
	if r.Temperature <= 25 {
		res.TempScore = 30
	} else {
		res.TempScore = max(0, 30-int((r.Temperature-25)*2))
		res.Feedback = append(res.Feedback, "Suhu terlalu tinggi untuk Cold Brew. Idealnya menggunakan air suhu ruang atau air dingin.")
	}

	// Time (BrewTime in seconds, so 12h = 43200s, 24h = 86400s)
	if r.BrewTime >= 43200 && r.BrewTime <= 86400 {
		res.TimeScore = 30
	} else {
		res.TimeScore = max(0, 30-int(math.Abs(float64(r.BrewTime)-64800)*0.0001))
		if r.BrewTime < 43200 {
			res.Feedback = append(res.Feedback, "Waktu ekstraksi Cold Brew biasanya minimal 12 jam.")
		}
	}

	// Grind Size Check
	if r.GrindSize != "coarse" {
		res.TotalScore -= 5
		res.Feedback = append(res.Feedback, "Cold Brew sangat disarankan menggunakan gilingan Coarse agar hasil lebih bersih.")
	}
}

func predictFlavor(r *Recipe, ratio float64) FlavorProfile {
	f := FlavorProfile{
		Acidity:    5,
		Sweetness:  5,
		Body:       5,
		Bitterness: 5,
	}

	// Temperature effects
	if r.Temperature > 95 {
		f.Bitterness += 3
		f.Acidity -= 2
	} else if r.Temperature < 88 && r.Method != "coldbrew" {
		f.Acidity += 3
		f.Sweetness -= 1
	}

	// Ratio effects
	if r.Method == "espresso" {
		f.Body += 3
		if ratio < 2.0 {
			f.Body += 2
			f.Sweetness += 1
		}
	} else {
		if ratio < 14.0 {
			f.Body += 2
		} else if ratio > 17.0 {
			f.Body -= 2
			f.Bitterness += 1
		}
	}

	// Time effects
	targetTime := 180
	if r.Method == "espresso" {
		targetTime = 30
	} else if r.Method == "frenchpress" {
		targetTime = 270
	} else if r.Method == "coldbrew" {
		targetTime = 64800
	}

	timeDiff := float64(r.BrewTime - targetTime)
	if timeDiff > 30 {
		f.Bitterness += 2
		f.Sweetness -= 1
	} else if timeDiff < -30 {
		f.Acidity += 2
		f.Body -= 1
	}

	// Clamp values 1-10
	clamp := func(v int) int {
		if v < 1 {
			return 1
		}
		if v > 10 {
			return 10
		}
		return v
	}

	f.Acidity = clamp(f.Acidity)
	f.Sweetness = clamp(f.Sweetness)
	f.Body = clamp(f.Body)
	f.Bitterness = clamp(f.Bitterness)

	return f
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
