package testutil

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func CmpOptions() []cmp.Option {
	options := []cmp.Option{
		cmp.Comparer(compareDecimals),
		cmp.Comparer(compareTime),
	}

	return options
}

func compareDecimals(a, b decimal.Decimal) bool {
	return a.Equal(b)
}

func compareTime(a, b time.Time) bool {
	return a.Equal(b)
}
