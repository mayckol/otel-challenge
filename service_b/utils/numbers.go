package utils

import "math"

func RoundToDecimal(value float64, decimalPoints int) float64 {
	shift := math.Pow(10, float64(decimalPoints))
	return math.Round(value*shift) / shift
}
