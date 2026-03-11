package utilities

import (
	"math/big"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func StringToDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

func Int64ToDecimal(amount int64, precision int32) (decimal.Decimal, error) {
	return decimal.NewFromBigInt(big.NewInt(amount), precision), nil
}
