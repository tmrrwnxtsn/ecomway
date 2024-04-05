package convert

import "math"

func BaseToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func CentsToBase(cents int64) float64 {
	return float64(cents) / 100
}
