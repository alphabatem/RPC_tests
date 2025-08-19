package cmd

import (
	"github.com/spf13/cobra"
)

// getAccountInfoCmd represents the getAccountInfo command
var getAccountInfoCmd = &cobra.Command{
	Use:   "getAccountInfo",
	Short: "Run performance tests for getAccountInfo RPC method",
	Long: `Run stress tests against Solana RPC endpoints using the getAccountInfo method.

This method tests individual account information retrieval performance. When multiple accounts 
are provided, the tool will rotate through them to distribute load evenly across different 
accounts and avoid overwhelming any single account.

Features:
• Account Rotation: Cycles through provided accounts for load distribution
• Real-time Progress: Visual progress bars with completion percentage and live statistics
• Comprehensive Metrics: Success rate, RPS, and latency statistics with dynamic unit formatting
• Flexible Input: Support for individual accounts or account files

Examples:
  # Test with a single account
  rpc_test getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --concurrency 5 --duration 30

  # Test with multiple accounts (will rotate between them)
  rpc_test getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --account vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg --concurrency 10 --duration 45

  # Test with accounts from a file (recommended for large lists)
  rpc_test getAccountInfo --account-file ./accounts.txt --concurrency 20 --duration 60 --limit 100`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMethodTest("getAccountInfo")
	},
}

func init() {
	RootCmd.AddCommand(getAccountInfoCmd)
}
