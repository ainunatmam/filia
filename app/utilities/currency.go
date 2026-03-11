package utilities

import (
	"errors"

	"github.com/shopspring/decimal"
)

func ToMinorUnit(amount decimal.Decimal, precision int32) (int64, error) {
    rounded := amount.Round(precision)

    multiplier := decimal.NewFromInt(1).Shift(precision)
    minor := rounded.Mul(multiplier)

    if !minor.IsInteger() {
        return 0, errors.New("invalid precision")
    }

    return minor.IntPart(), nil
}

func ToMajorString(amount uint64, precision int32) string {
    major := decimal.NewFromUint64(amount).Shift(-precision)
    return major.String()
}