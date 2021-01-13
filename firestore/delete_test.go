package firestore_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/blaggotech/gostcalc/firestore"
)

func Test_DailyDeleteCalculator_Calculate(t *testing.T) {
	calc := &firestore.DailyDeleteCalculator{}
	want := 0.02

	dailyDeletes := big.NewInt(100000)
	res, err := calc.Calculate(context.Background(), dailyDeletes)
	if err != nil {
		t.Fatalf("unable to calculate daily deletes: %v", err)
	}

	got, _ := res.Float64()

	if want != got {
		t.Errorf("want Calculate() result to be %v, got %v", want, got)
	}
}

func Test_MonthlyDeleteCalculator_Calculate(t *testing.T) {
	calc := &firestore.MonthlyDeleteCalculator{
		D: &firestore.DailyDeleteCalculator{},
	}
	// 0.02 * 30
	want := 0.6 

	dailyDeletes := big.NewInt(100000)
	res, err := calc.Calculate(context.Background(), dailyDeletes)
	if err != nil {
		t.Fatalf("unable to calculate daily deletes: %v", err)
	}

	got, _ := res.Float64()

	if want != got {
		t.Errorf("want Calculate() result to be %v, got %v", want, got)
	}
}
