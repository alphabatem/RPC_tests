package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "rpc_test",
	Short: "A tool for testing Solana RPC endpoints",
	Long: `A CLI tool for testing and benchmarking Solana RPC endpoints with customizable concurrency and account inputs.

Example:
  # Run a test using getAccountInfo method
  rpc_test getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --url https://api.mainnet-beta.solana.com

  # Run a test using getProgramAccounts method
  rpc_test getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --concurrency 10 --duration 30

  # Seed account data for a program
  rpc_test seed --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --output ./accounts`,
}

func init() {
	// Common flags for all commands
	RootCmd.PersistentFlags().StringVarP(&rpcURL, "url", "u", "https://api.mainnet-beta.solana.com", "RPC endpoint URL")
	RootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent requests")
	RootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 10, "Test duration in seconds")
	RootCmd.PersistentFlags().StringArrayVarP(&accounts, "account", "a", []string{}, "Account addresses to use in tests (can be specified multiple times)")
	RootCmd.PersistentFlags().StringVarP(&accountsFile, "account-file", "f", "", "File containing account addresses (one per line)")
	RootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "Limit the number of accounts/programs to process (0 for no limit)")
}

// Execute adds all child commands to the root command and executes it
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
