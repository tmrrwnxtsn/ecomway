package convert

import (
	"math"
	"strconv"
)

func BaseToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func CentsToBase(cents int64) float64 {
	return float64(cents) / 100
}

func FloatWithoutTrailingZeroes(input float64, precision int) string {
	// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
	// https://stackoverflow.com/questions/31289409/format-a-float-to-n-decimal-places-and-no-trailing-zeros
	output := math.Pow(10, float64(precision))
	multiplexed := input * output
	roundedPrecisionFloat := float64(int(multiplexed+math.Copysign(0.5, multiplexed))) / output

	return strconv.FormatFloat(roundedPrecisionFloat, 'f', precision, 64)
}
