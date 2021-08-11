package firestore_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/royge/gostcalc/firestore"
)

func Test_DailyReadCalculator_Calculate(t *testing.T) {
	calc := &firestore.DailyReadCalculator{}
	want := 0.21

	dailyReads := big.NewInt(400000)
	res, err := calc.Calculate(context.Background(), dailyReads)
	if err != nil {
		t.Fatalf("unable to calculate daily reads: %v", err)
	}

	got, _ := res.Float64()

	if want != got {
		t.Errorf("want Calculate() result to be %v, got %v", want, got)
	}
}

func Test_MonthlyReadCalculator_Calculate(t *testing.T) {
	calc := &firestore.MonthlyReadCalculator{
		D: &firestore.DailyReadCalculator{},
	}
	// 0.21 * 30
	want := 6.3

	dailyReads := big.NewInt(400000)
	res, err := calc.Calculate(context.Background(), dailyReads)
	if err != nil {
		t.Fatalf("unable to calculate daily reads: %v", err)
	}

	got, _ := res.Float64()

	if want != got {
		t.Errorf("want Calculate() result to be %v, got %v", want, got)
	}
}
