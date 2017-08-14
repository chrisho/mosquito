package helper

import (
	"math"
	"fmt"
	"strconv"
)

// val:值，place:精度（多少个小数位）
func PhpRound(val float64, places int) float64 {
	var t float64
	f := math.Pow10(places)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

// val:值，place:精度（多少个小数位）
func GolangRound(val float64, places int) float64 {
	format := `%.` + strconv.Itoa(places) + `f`
	floatStr := fmt.Sprintf(format, val)
	result, _ := strconv.ParseFloat(floatStr, 10)

	return result
}
