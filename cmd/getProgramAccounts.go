package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	programs     []string
	programsFile string
)

// getProgramAccountsCmd represents the getProgramAccounts command
var getProgramAccountsCmd = &cobra.Command{
	Use:   "getProgramAccounts",
	Short: "Run performance tests for getProgramAccounts RPC method",
	Long: `Run stress tests against Solana RPC endpoints using the getProgramAccounts method.

This method tests program account enumeration performance, which is typically one of the most 
resource-intensive RPC operations. It fetches all accounts owned by specified programs, making 
it ideal for testing RPC endpoint capabilities under heavy load.

Features:
• Program Rotation: Cycles through provided programs for load distribution
• Real-time Progress: Visual progress bars with completion percentage and live statistics
• Comprehensive Metrics: Success rate, RPS, and latency statistics with dynamic unit formatting
• Flexible Input: Support for individual programs or program files
• Resource Intensive: Tests the most demanding RPC operation for comprehensive benchmarking

Note: Use --program flag (not --account) and --program-file (not --account-file) for this command.

Examples:
  # Test with a single program
  rpc_test getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --concurrency 5 --duration 30

  # Test with multiple programs (will rotate between them)
  rpc_test getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --program 2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c --concurrency 10 --duration 45

  # Test with programs from a file (recommended for multiple programs)
  rpc_test getProgramAccounts --program-file ./programs.txt --concurrency 20 --duration 60 --limit 10`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load programs from file if provided
		if programsFile != "" {
			data, err := os.ReadFile(programsFile)
			if err != nil {
				log.Fatalf("Failed to read programs file: %v", err)
			}
			// Parse programs from file
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					programs = append(programs, line)
				}
			}
		}

		if len(programs) == 0 {
			log.Fatalf("No programs provided. Use --program or --program-file to specify programs")
		}

		// Apply limit if specified
		totalPrograms := len(programs)
		if limit > 0 && limit < totalPrograms {
			programs = programs[:limit]
			log.Printf("Limiting to %d programs out of %d available", limit, totalPrograms)
		}

		// Use programs as accounts for the underlying test runner
		accounts = programs

		RunMethodTest("getProgramAccounts")
	},
}

func init() {
	RootCmd.AddCommand(getProgramAccountsCmd)

	// Add program-specific flags
	getProgramAccountsCmd.Flags().StringArrayVarP(&programs, "program", "p", []string{}, "Program addresses to use in tests (can be specified multiple times)")
	getProgramAccountsCmd.Flags().StringVarP(&programsFile, "program-file", "f", "", "File containing program addresses (one per line)")

	// Override the account-file flag to avoid confusion
	getProgramAccountsCmd.Flags().StringVarP(&accountsFile, "account-file", "", "", "")
	getProgramAccountsCmd.Flags().MarkHidden("account-file")

	// Override the account flag to avoid confusion
	getProgramAccountsCmd.Flags().StringArrayVarP(&accounts, "account", "", []string{}, "")
	getProgramAccountsCmd.Flags().MarkHidden("account")
}
