package firestore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/blaggotech/gostcalc/firestore"
	"github.com/google/uuid"
)

func Test_MonthlyNetworkingCalculator_Calculate(t *testing.T) {
	// Cebu city population as of 2015.
	// cebuPopulation := int64(922611)
	cebuPopulation := int64(1000000)

	// Number of QR transactions per day.
	numOfTxn := int64(10)

	data := map[string]interface{}{
		"id":          uuid.New(),
		"profile_id":  uuid.New(),
		"merchant_id": uuid.New(),
	}

	doc, err := json.Marshal(&data)
	if err != nil {
		t.Errorf("unable to marshal data: %v", err)
	}

	calc := &firestore.MonthlyNetworkingCalculator{
		D: &firestore.DailyNetworkingCalculator{
			Document: doc,
		},
		Price: firestore.IngressPricePerGB,
	}

	cost, err := calc.Calculate(
		context.Background(),
		big.NewInt(int64(cebuPopulation)*numOfTxn),
	)
	if err != nil {
		t.Fatalf("unable to calculate ingress cost: %v", err)
	}

	fmt.Println("Networking Ingress Cost: $", cost)
}
