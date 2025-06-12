package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// localRpcCmd represents the localRpc command
var localRpcCmd = &cobra.Command{
	Use:   "localRpc",
	Short: "Test against a local RPC endpoint at localhost:8080",
	Long: `Run performance tests against a local Solana RPC endpoint running on localhost:8080.

Example:
  # Test getAccountInfo against local RPC
  rpc_test localRpc getAccountInfo --account 7Xnw7aDxJu1CxPPEkz9ttfGSn2bpH3R1GYYziJxTCv3e --concurrency 5 --duration 30

  # Test getProgramAccounts against local RPC
  rpc_test localRpc getProgramAccounts --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --concurrency 10

  # Seed accounts from a local program
  rpc_test localRpc seed --program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA --output ./local_accounts`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Override the RPC URL to use localhost:8080
		rpcURL = "http://localhost:8080"
		fmt.Printf("Using local RPC endpoint: %s\n", rpcURL)
	},
}

func init() {
	RootCmd.AddCommand(localRpcCmd)

	// Add subcommands
	localRpcCmd.AddCommand(getAccountInfoCmd)
	localRpcCmd.AddCommand(getMultipleAccountsCmd)
	localRpcCmd.AddCommand(getProgramAccountsCmd)
	localRpcCmd.AddCommand(seedCmd)
}
