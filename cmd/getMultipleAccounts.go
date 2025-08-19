package cmd

import (
	"github.com/spf13/cobra"
)

// getMultipleAccountsCmd represents the getMultipleAccounts command
var getMultipleAccountsCmd = &cobra.Command{
	Use:   "getMultipleAccounts",
	Short: "Run performance tests for getMultipleAccounts RPC method",
	Long: `Run stress tests against Solana RPC endpoints using the getMultipleAccounts method.

This method automatically batches accounts from your provided list into groups of 5-15 accounts 
per request (randomized) to simulate realistic batch retrieval patterns. Each worker thread 
will rotate through your account list and create different batch combinations.

Features:
• Automatic Batching: Groups 5-15 accounts per request (randomized for variety)
• Account Rotation: Cycles through provided accounts for load distribution
• Real-time Progress: Visual progress bars with completion percentage and live statistics
• Comprehensive Metrics: Success rate, RPS, and latency statistics with dynamic unit formatting

Examples:
  # Test with multiple accounts (will be batched automatically)
  rpc_test getMultipleAccounts --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --account vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg --concurrency 5 --duration 30

  # Test with accounts from a file (recommended for large lists)
  rpc_test getMultipleAccounts --account-file ./accounts.txt --concurrency 10 --duration 60

  # Test with custom settings and account limit
  rpc_test getMultipleAccounts --account-file ./accounts.txt --limit 100 --concurrency 15 --duration 45`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMethodTest("getMultipleAccounts")
	},
}

func init() {
	RootCmd.AddCommand(getMultipleAccountsCmd)
}
