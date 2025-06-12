package cmd

import (
	"github.com/spf13/cobra"
)

// getMultipleAccountsCmd represents the getMultipleAccounts command
var getMultipleAccountsCmd = &cobra.Command{
	Use:   "getMultipleAccounts",
	Short: "Run performance tests for getMultipleAccounts RPC method",
	Long: `Run stress tests against Solana RPC endpoints using the getMultipleAccounts method.

Example:
  # Test with a single group of accounts
  rpc_test getMultipleAccounts --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e,vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg --concurrency 5 --duration 30

  # Test with multiple account groups
  rpc_test getMultipleAccounts --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --account vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg --concurrency 10

  # Test with accounts from a file
  rpc_test getMultipleAccounts --account-file ./accounts.txt --url https://api.mainnet-beta.solana.com --concurrency 20 --duration 60`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMethodTest("getMultipleAccounts")
	},
}

func init() {
	RootCmd.AddCommand(getMultipleAccountsCmd)
}
