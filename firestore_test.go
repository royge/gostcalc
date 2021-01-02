package main

import (
	"context"
	"testing"
)

func Test_WriteCalculator_Calculate(t *testing.T) {
	calc := &WriteCalculator{}

	_, err := calc.Calculate(context.Background(), numberOfWrites)
	if err != nil {
		t.Errorf("unable to calculate cost: %v", err)
	}
}
