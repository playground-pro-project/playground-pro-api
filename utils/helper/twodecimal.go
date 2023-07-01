package helper

import "math"

func TwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
