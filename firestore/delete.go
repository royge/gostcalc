package firestore

import (
	"context"
	"math/big"
)

const (
	DeleteUnitPrice  = 0.02
	FreeDeletesDaily = 20000
)

type DailyDeleteCalculator struct{}

func (dw *DailyDeleteCalculator) Calculate(_ context.Context, count *big.Int) (*big.Float, error) {
	free := big.NewInt(FreeDeletesDaily)
	b := count.Sub(count, free)
	billable := new(big.Float).SetInt(b)

	unit := new(big.Float).SetInt64(Unit)
	unitPrice := big.NewFloat(DeleteUnitPrice)

	q := billable.Quo(billable, unit)

	daily := q.Mul(q, unitPrice)

	return daily, nil
}

type MonthlyDeleteCalculator struct {
	D *DailyDeleteCalculator
}

func (mw *MonthlyDeleteCalculator) Calculate(ctx context.Context, count *big.Int) (*big.Float, error) {
	daily, err := mw.D.Calculate(ctx, count)
	if err != nil {
		return nil, err
	}

	days := new(big.Float).SetInt64(int64(MonthNumOfDays))
	monthly := daily.Mul(daily, days)

	return monthly, nil
}
