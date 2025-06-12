package cmd

import (

	"github.com/spf13/cobra"
)

// getAccountInfoCmd represents the getAccountInfo command
var getAccountInfoCmd = &cobra.Command{
	Use:   "getAccountInfo",
	Short: "Run performance tests for getAccountInfo RPC method",
	Long: `Run stress tests against Solana RPC endpoints using the getAccountInfo method.

Example:
  # Test with a single account
  rpc_test getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --concurrency 5 --duration 30

  # Test with multiple accounts using command line
  rpc_test getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --account vines1vzrYbzLMRdu58ou5XTby4qAqVRLmqo36NKPTg --concurrency 10

  # Test with accounts from a file
  rpc_test getAccountInfo --account-file ./accounts.txt --url https://api.mainnet-beta.solana.com --concurrency 20 --duration 60`,
	Run: func(cmd *cobra.Command, args []string) {
		RunMethodTest("getAccountInfo")
	},
}

func init() {
	RootCmd.AddCommand(getAccountInfoCmd)
}
