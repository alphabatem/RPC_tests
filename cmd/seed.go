package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"rpc_test/methods"

	"github.com/spf13/cobra"
)

var (
	outputFile string
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed account addresses from programs for testing purposes",
	Long: `Fetch program accounts and save their addresses to a file for testing purposes.

This command fetches all accounts owned by specified programs and saves their addresses 
to a text file. This is essential for preparing test data before running performance tests 
with getAccountInfo or getMultipleAccounts methods.

Features:
• Bulk Account Fetching: Retrieves all accounts from specified programs efficiently
• File Output: Saves account addresses in a clean text format (one per line)
• Limit Support: Control the number of accounts to fetch with --limit flag
• Directory Creation: Automatically creates output directories if they don't exist
• Multiple Programs: Support for fetching from multiple programs simultaneously

Use Cases:
• Prepare account lists for performance testing
• Generate test data for CI/CD pipelines
• Create account datasets for load testing scenarios
• Build custom account collections for specific testing needs

Examples:
  # Seed accounts for a single program
  rpc_test seed --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --output accounts.txt

  # Seed accounts with a limit (recommended for testing)
  rpc_test seed --program 2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c --output accounts.txt --limit 1000

  # Seed from multiple programs 
  rpc_test seed --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --program 2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c --output ./data/accounts.txt

  # Seed from programs listed in a file
  rpc_test seed --program-file ./programs.txt --output ./data/test_accounts.txt --limit 500`,
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

		// Create output directory if needed
		outputDir := filepath.Dir(outputFile)
		if outputDir != "." && outputDir != "" {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				log.Fatalf("Failed to create output directory: %v", err)
			}
		}

		fmt.Printf("Fetching accounts for %d programs\n", len(programs))

		for _, program := range programs {
			fmt.Printf("Processing program: %s\n", program)
			err := seedProgramAccounts(program, outputFile)
			if err != nil {
				log.Printf("Error processing program %s: %v", program, err)
			}
		}
	},
}

// seedProgramAccounts fetches and saves program accounts
func seedProgramAccounts(programAddress string, outputFile string) error {
	// Create RPC client
	rpcTest := methods.NewRPCTest(rpcURL, apiKey)

	// Seed program accounts
	return rpcTest.SeedProgramAccounts(programAddress, outputFile, limit)
}

func init() {
	RootCmd.AddCommand(seedCmd)

	// Add program-specific flags
	seedCmd.Flags().StringArrayVarP(&programs, "program", "p", []string{}, "Program addresses to fetch accounts for (can be specified multiple times)")
	seedCmd.Flags().StringVarP(&programsFile, "program-file", "f", "", "File containing program addresses (one per line)")
	seedCmd.Flags().StringVarP(&outputFile, "output", "o", "accounts.txt", "Output file to store account addresses")

	// Override the account-file flag to avoid confusion
	seedCmd.Flags().StringVarP(&accountsFile, "account-file", "", "", "")
	seedCmd.Flags().MarkHidden("account-file")

	// Override the account flag to avoid confusion
	seedCmd.Flags().StringArrayVarP(&accounts, "account", "", []string{}, "")
	seedCmd.Flags().MarkHidden("account")
}
