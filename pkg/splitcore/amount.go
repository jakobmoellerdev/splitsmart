package splitcore

import "math"

type Amount float64

func (a Amount) Equals(b Amount) bool {
	return Round(a) == Round(b)
}

var RoundingDigitCount = 4

func Round(amount Amount) Amount {
	return Amount(round(float64(amount), RoundingDigitCount))
}

// round Rounds to nearest like 12.3456 -> 12.35.
func round(val float64, precision int) float64 {
	return math.Round(val*(math.Pow10(precision))) / math.Pow10(precision)
}
