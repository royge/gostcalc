package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/blaggotech/gostcalc/firestore"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	dailyTxn   int64
	population int64
)

// RegisterFirestore register/initialize CLI command to calculate firestore
// costs.
func RegisterFirestore() {
	rootCmd.AddCommand(networkCmd)
	rootCmd.AddCommand(storageCmd)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(readCmd)

	networkCmd.Flags().Int64VarP(
		&dailyTxn,
		"count",
		"c",
		10,
		"Total number of daily transactions",
	)

	networkCmd.Flags().Int64VarP(
		&population,
		"population",
		"p",
		1000000,
		"Total number of active users",
	)

	storageCmd.Flags().Int64VarP(
		&dailyTxn,
		"count",
		"c",
		10,
		"Total number of daily transactions",
	)

	storageCmd.Flags().Int64VarP(
		&population,
		"population",
		"p",
		1000000,
		"Total number of active users",
	)

	writeCmd.Flags().Int64VarP(
		&dailyTxn,
		"count",
		"c",
		10,
		"Total number of daily transactions",
	)

	writeCmd.Flags().Int64VarP(
		&population,
		"population",
		"p",
		1000000,
		"Total number of active users",
	)

	deleteCmd.Flags().Int64VarP(
		&dailyTxn,
		"count",
		"c",
		10,
		"Total number of daily transactions",
	)

	deleteCmd.Flags().Int64VarP(
		&population,
		"population",
		"p",
		1000000,
		"Total number of active users",
	)

	readCmd.Flags().Int64VarP(
		&dailyTxn,
		"count",
		"c",
		10,
		"Total number of daily transactions",
	)

	readCmd.Flags().Int64VarP(
		&population,
		"population",
		"p",
		1000000,
		"Total number of active users",
	)
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Calculate network ingress costs.",
	Long:  "Calculate network ingress costs.",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]interface{}{
			"id":          uuid.New(),
			"profile_id":  uuid.New(),
			"merchant_id": uuid.New(),
		}

		doc, err := json.Marshal(&data)
		if err != nil {
			log.Fatalf("unable to marshal data: %v", err)
			os.Exit(1)
		}

		calc := &firestore.MonthlyNetworkingCalculator{
			D: &firestore.DailyNetworkingCalculator{
				Document: doc,
			},
			Price: firestore.IngressPricePerGB,
		}

		cost, err := calc.Calculate(
			context.Background(),
			big.NewInt(population*dailyTxn),
		)
		if err != nil {
			log.Fatalf("unable to calculate ingress cost: %v", err)
			os.Exit(1)
		}

		fmt.Println("Estimated Networking Cost: $", cost)
	},
}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Calculate firestore storage costs.",
	Long:  "Calculate firestore storage costs.",
	Run: func(cmd *cobra.Command, args []string) {
		calc := &firestore.MonthlyStorageCalculator{
			D: &firestore.DailyStorageCalculator{
				Document: &firestore.Document{
					ID: uuid.New().String(),
					Collection: fmt.Sprintf(
						"prod-qr/%s/qr-records",
						uuid.New().String(),
					),
					Data: map[string]interface{}{
						"merchant_id":     uuid.New().String(),
						"merchant_qr_id":  uuid.New().String(),
						"profile_qr_id":   uuid.New().String(),
						"date_created":    time.Now(),
						"type":            1,
						// "is_auto_scanout": false,
						"is_auto_scanout": map[string]interface{}{
							"Bool": false,
							"Valid": false,
						},
					},
					SingleFieldIndexes: []map[string]interface{}{
						{
							"date_created": time.Now(),
						},
					},
					CompositeIndexes: []map[string]interface{}{
						{
							"merchant_id":  uuid.New().String(),
							"date_created": time.Now(),
						},
						{
							"merchant_id":  uuid.New().String(),
							"type":         1,
							"date_created": time.Now(),
						},
						{
							"type":         1,
							"date_created": time.Now(),
						},
					},
				},
			},
			Price: firestore.PricePerGB,
		}

		cost, err := calc.Calculate(
			context.Background(),
			big.NewInt(population*dailyTxn),
		)
		if err != nil {
			log.Fatalf("unable to calculate storage cost: %v", err)
			os.Exit(1)
		}

		fmt.Println("Estimated Storage Cost: $", cost)
	},
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Calculate firestore write costs.",
	Long:  "Calculate firestore write costs.",
	Run: func(cmd *cobra.Command, args []string) {
		calc := &firestore.MonthlyWriteCalculator{
			D: &firestore.DailyWriteCalculator{},
		}

		dailyWrites := big.NewInt(population * dailyTxn)
		cost, err := calc.Calculate(
			context.Background(),
			dailyWrites,
		)
		if err != nil {
			log.Fatalf("unable to calculate daily writes: %v", err)
		}

		fmt.Println("Estimated Writes Cost: $", cost)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Calculate firestore delete costs.",
	Long:  "Calculate firestore delete costs.",
	Run: func(cmd *cobra.Command, args []string) {
		calc := &firestore.MonthlyDeleteCalculator{
			D: &firestore.DailyDeleteCalculator{},
		}

		dailyDeletes := big.NewInt(population * dailyTxn)
		cost, err := calc.Calculate(
			context.Background(),
			dailyDeletes,
		)
		if err != nil {
			log.Fatalf("unable to calculate daily deletes: %v", err)
		}

		fmt.Println("Estimated Deletes Cost: $", cost)
	},
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Calculate firestore read costs.",
	Long:  "Calculate firestore read costs.",
	Run: func(cmd *cobra.Command, args []string) {
		calc := &firestore.MonthlyReadCalculator{
			D: &firestore.DailyReadCalculator{},
		}

		dailyReads := big.NewInt(population * dailyTxn)
		cost, err := calc.Calculate(
			context.Background(),
			dailyReads,
		)
		if err != nil {
			log.Fatalf("unable to calculate daily reads: %v", err)
		}

		fmt.Println("Estimated Reads Cost: $", cost)
	},
}
