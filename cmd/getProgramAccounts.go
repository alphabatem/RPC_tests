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

Example:
  # Test with a single program
  rpc_test getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --concurrency 5 --duration 30

  # Test with multiple programs
  rpc_test getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --program 9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin --concurrency 10

  # Test with programs from a file
  rpc_test getProgramAccounts --program-file ./programs.txt --url https://api.mainnet-beta.solana.com --concurrency 20 --duration 60`,
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
