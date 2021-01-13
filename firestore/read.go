package firestore

import (
	"context"
	"math/big"
)

const (
	ReadUnitPrice  = 0.06
	FreeReadsDaily = 50000
)

type DailyReadCalculator struct{}

func (dw *DailyReadCalculator) Calculate(_ context.Context, count *big.Int) (*big.Float, error) {
	free := big.NewInt(FreeReadsDaily)
	b := count.Sub(count, free)
	billable := new(big.Float).SetInt(b)

	unit := new(big.Float).SetInt64(Unit)
	unitPrice := big.NewFloat(ReadUnitPrice)

	q := billable.Quo(billable, unit)

	daily := q.Mul(q, unitPrice)

	return daily, nil
}

type MonthlyReadCalculator struct {
	D *DailyReadCalculator
}

func (mw *MonthlyReadCalculator) Calculate(ctx context.Context, count *big.Int) (*big.Float, error) {
	daily, err := mw.D.Calculate(ctx, count)
	if err != nil {
		return nil, err
	}

	days := new(big.Float).SetInt64(int64(MonthNumOfDays))
	monthly := daily.Mul(daily, days)

	return monthly, nil
}
