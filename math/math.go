package math

import (
	"math"

	"github.com/shopspring/decimal"
)

// Round (5.555, 2) -> 5.56
func Round(v float64, precision int32) float64 {
	dValue := decimal.NewFromFloat(v)
	out, _ := dValue.Round(precision).Float64()
	return out
}

// DecimalFloor (5.555, 2) -> 5.55
func DecimalFloor(v decimal.Decimal, precision int32) decimal.Decimal {
	dMulti := decimal.NewFromFloat(math.Pow10(int(precision)))
	return v.Mul(dMulti).Floor().Div(dMulti)
}

// Floor (5.555, 2) -> 5.55
func Floor(v float64, precision int32) float64 {
	out, _ := DecimalFloor(NewFromFloat(v), precision).Float64()
	return out
}

// DecimalCeil (5.551, 2) -> 5.56
func DecimalCeil(v decimal.Decimal, precision int32) decimal.Decimal {
	dMulti := decimal.NewFromFloat(math.Pow10(int(precision)))
	return v.Mul(dMulti).Ceil().Div(dMulti)
}

// Ceil (5.551, 2) -> 5.56
func Ceil(v float64, precision int32) float64 {
	out, _ := DecimalCeil(NewFromFloat(v), precision).Float64()
	return out
}

// NewFromFloat max 12 digits decimal
func NewFromFloat(v float64) decimal.Decimal {
	return decimal.NewFromFloatWithExponent(v, -12)
}

// FloatToInt (15.6666, 4) -> 156666
func FloatToInt(f float64, precision int) int {
	dMulti := decimal.NewFromFloat(math.Pow10(precision))
	dF := decimal.NewFromFloat(f)
	return int(dF.Mul(dMulti).Round(0).IntPart())
}

// FloatFromInt (156666, 4) -> 15.6666
func FloatFromInt(i int, precision int) float64 {
	dMulti := decimal.NewFromFloat(math.Pow10(precision))
	dI := decimal.NewFromFloat(float64(i))
	out, _ := dI.Div(dMulti).Round(int32(precision)).Float64()
	return out
}
