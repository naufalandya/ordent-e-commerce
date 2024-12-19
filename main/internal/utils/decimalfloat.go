package utils

import (
	"github.com/shopspring/decimal"
)

func DecimalFromFloat64(value float64) decimal.Decimal {
	return decimal.NewFromFloat(value)
}
