package firestore

import (
	"context"
	"math/big"
)

const (
	// MonthlyFreeIngress is the 10GB free monthly ingress.
	MonthlyFreeIngress = 10737418240

	// IngressPricePerGB is ingress price per GB after free.
	IngressPricePerGB = 0.12
)

type DailyNetworkingCalculator struct {
	// Document in transit.
	Document []byte
}

func (dn *DailyNetworkingCalculator) Calculate(_ context.Context, count *big.Int) (*big.Float, error) {
	// Get document size in bytes.
	size := big.NewInt(int64(len(dn.Document)))

	daily := size.Mul(size, count)

	return new(big.Float).SetInt(daily), nil
}

type MonthlyNetworkingCalculator struct {
	D *DailyNetworkingCalculator

	// Unit Price.
	// Price per GB.
	Price float64
}

func (mn *MonthlyNetworkingCalculator) Calculate(ctx context.Context, count *big.Int) (*big.Float, error) {
	daily, err := mn.D.Calculate(ctx, count)
	if err != nil {
		return new(big.Float), err
	}

	days := new(big.Float).SetInt64(MonthNumOfDays)
	monthly := daily.Mul(daily, days)

	free := new(big.Float).SetInt64(MonthlyFreeIngress)

	monthly = monthly.Sub(monthly, free)
	monthly = monthly.Quo(monthly, new(big.Float).SetInt64(OneGB))

	cost := monthly.Mul(monthly, new(big.Float).SetFloat64(mn.Price))

	return cost, nil
}
