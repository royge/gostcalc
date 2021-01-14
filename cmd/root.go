package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gostcalc",
	Short: "Google cloud cost calculator.",
	Long:  "This tool will try to calculate the cost estimations based on input data.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Welcome! try gostcalc -h for usage instructions.")
	},
}

// Execute commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Register calls every func in funcs as command registration.
func Register(funcs ...func()) {
	for _, fn := range funcs {
		fn()
	}
}
