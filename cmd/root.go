package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var RootCmd = &cobra.Command{
	Use:   "rpc_test",
	Short: "A comprehensive CLI tool for stress testing and benchmarking Solana RPC endpoints",
	Long: `A comprehensive CLI tool for stress testing and benchmarking Solana RPC endpoints with customizable concurrency and account inputs.

Features:
• Comprehensive Test Suite: Run all RPC methods simultaneously with 'runall' command
• Dynamic Configuration: Generate and load test configurations with API keys  
• Smart Progress Tracking: Real-time progress bars and detailed statistics
• Dynamic Latency Display: Automatic unit selection (μs, ms, s) based on performance
• Dual RPC Architecture: Separate seeding vs testing RPCs for advanced workflows
• Account Management: Seed accounts from programs or use custom account lists
• Performance Metrics: Comprehensive latency, RPS, and success rate statistics

Supported RPC Methods:
• getAccountInfo: Test account information retrieval with account rotation
• getMultipleAccounts: Test batch account retrieval (5-15 accounts per request)
• getProgramAccounts: Test program account enumeration

Examples:
  # Run comprehensive test suite (recommended)
  rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com

  # Test individual methods
  rpc_test getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --concurrency 10 --duration 30
  rpc_test getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --concurrency 5 --duration 15

  # Seed account data for testing
  rpc_test seed --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --output ./accounts.txt --limit 1000`,
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
