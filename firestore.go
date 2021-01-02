package main

import (
	"context"
	"errors"
)

// WriteCalculator defines how to calculate Firestore write costs.
type WriteCalculator struct{
	document []byte
}

func (w *WriteCalculator) Calculate(_ context.Context, count int64) (float64, error) {
	return 0, errors.New("not yet implemented")
}
