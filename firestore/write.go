package firestore

import (
	"context"
	"math/big"
)

const (
	Unit            = 100000
	WriteUnitPrice  = 0.18
	FreeWritesDaily = 20000
)

type DailyWriteCalculator struct{}

func (dw *DailyWriteCalculator) Calculate(_ context.Context, count *big.Int) (*big.Float, error) {
	free := big.NewInt(FreeWritesDaily)
	b := count.Sub(count, free)
	billable := new(big.Float).SetInt(b)

	unit := new(big.Float).SetInt64(Unit)
	unitPrice := big.NewFloat(WriteUnitPrice)

	q := billable.Quo(billable, unit)

	daily := q.Mul(q, unitPrice)

	return daily, nil
}

type MonthlyWriteCalculator struct {
	D *DailyWriteCalculator
}

func (mw *MonthlyWriteCalculator) Calculate(ctx context.Context, count *big.Int) (*big.Float, error) {
	daily, err := mw.D.Calculate(ctx, count)
	if err != nil {
		return nil, err
	}

	days := new(big.Float).SetInt64(int64(MonthNumOfDays))
	monthly := daily.Mul(daily, days)

	return monthly, nil
}
