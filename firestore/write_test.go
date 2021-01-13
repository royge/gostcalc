package firestore_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/blaggotech/gostcalc/firestore"
)

func Test_DailyWriteCalculator_Calculate(t *testing.T) {
	calc := &firestore.DailyWriteCalculator{}
	want := 0.14

	dailyWrites := big.NewInt(100000)
	res, err := calc.Calculate(context.Background(), dailyWrites)
	if err != nil {
		t.Fatalf("unable to calculate daily writes: %v", err)
	}

	got, _ := res.Float64()

	if want != got {
		t.Errorf("want Calculate() result to be %v, got %v", want, got)
	}
}

func Test_MonthlyWriteCalculator_Calculate(t *testing.T) {
	calc := &firestore.MonthlyWriteCalculator{
		D: &firestore.DailyWriteCalculator{},
	}
	// 0.14 * 30
	want := 4.20

	dailyWrites := big.NewInt(100000)
	res, err := calc.Calculate(context.Background(), dailyWrites)
	if err != nil {
		t.Fatalf("unable to calculate daily writes: %v", err)
	}

	got, _ := res.Float64()

	if want != got {
		t.Errorf("want Calculate() result to be %v, got %v", want, got)
	}
}
