package firestore_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/blaggotech/gostcalc/firestore"
	"github.com/google/uuid"
)

func Test_MonthlyStorageCalculator_Calculate(t *testing.T) {
	// Cebu city population as of 2015.
	// cebuPopulation := int64(922611)
	cebuPopulation := int64(1000000)

	// Number of QR transactions per day.
	numOfTxn := int64(10)

	calc := &firestore.MonthlyStorageCalculator{
		D: &firestore.DailyStorageCalculator{
			Document: &firestore.Document{
				ID: uuid.New().String(),
				Collection: fmt.Sprintf(
					"profiles/%s/logs",
					uuid.New().String(),
				),
				Data: map[string]interface{}{
					"merchant_id": uuid.New().String(),
					"created":     time.Now(),
				},
				SingleFieldIndexes: []map[string]interface{}{
					{
						"created": time.Now(),
					},
				},
				CompositeIndexes: []map[string]interface{}{
					{
						"merchant_id": uuid.New().String(),
						"created":     time.Now(),
					},
				},
			},
		},
		Price: firestore.PricePerGB,
	}

	cost, err := calc.Calculate(
		context.Background(),
		big.NewInt(int64(cebuPopulation)*numOfTxn),
	)
	if err != nil {
		t.Fatalf("unable to calculate storage cost: %v", err)
	}

	fmt.Println("Storage Cost: $", cost)
}

func Test_Document_Size(t *testing.T) {
	doc := &firestore.Document{
		ID:         "my_task_id",
		Collection: "users/jeff/tasks",
		Data: map[string]interface{}{
			"type":        "Personal",
			"done":        false,
			"priority":    1,
			"description": "Learn Cloud Firestore",
			"created":     time.Now(),
		},
		CompositeIndexes: []map[string]interface{}{
			{
				"done":     false,
				"priority": 1,
			},
			// {
			// 	"type":    "Personal",
			// 	"created": time.Now(),
			// },
		},
	}

	// Breakdown of document size:
	//
	// Document name:
	// "users" = 5 + 1 = 6
	// "jeff" = 4 + 1 = 5
	// "tasks" = 5 + 1 = 6
	// "my_task_id" = 10 + 1 = 11
	// padding = 16
	//
	// Document data:
	// "type": "Personal" = (4 + 1 = 5) + (8 + 1 = 9) = 14
	// "done": false = (4 + 1 = 5) + 1 = 6
	// "priority": 1 = (8 + 1 = 9) + 8 = 17
	// "description": "Learn Cloud Firestore" = (11 + 1 = 12) + (21 + 1 = 22) = 34
	// "created": time.Now() = (7 + 1 = 8) + 8 = 16
	// padding = 32

	// Total
	// Document size: 163
	// Composite index size: 112
	want := int64(163 + 112) // bytes

	got := doc.Size()

	if got != want {
		t.Errorf("want Size() %v bytes, got %v bytes", want, got)
	}
}
