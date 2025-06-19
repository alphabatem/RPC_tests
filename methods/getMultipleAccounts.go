package methods

import (
	"context"
	"fmt"
	"strings"

	"github.com/gagliardetto/solana-go"
)

// GetMultipleAccounts fetches information for multiple accounts at once
func (r *RPCTest) GetMultipleAccounts(accountsStr ...string) error {

	// Parse the account addresses
	pubKeys := make([]solana.PublicKey, 0, len(accountsStr))
	for _, addrStr := range accountsStr {
		addrStr = strings.TrimSpace(addrStr)
		if addrStr == "" {
			continue
		}

		pubKey, err := solana.PublicKeyFromBase58(addrStr)
		if err != nil {
			return fmt.Errorf("invalid account address '%s': %v", addrStr, err)
		}
		pubKeys = append(pubKeys, pubKey)
	}

	if len(pubKeys) == 0 {
		return fmt.Errorf("no valid account addresses provided")
	}

	// Fetch multiple accounts
	_, err := r.rpc.GetMultipleAccounts(
		context.Background(),
		pubKeys...,
	)
	
	if err != nil {
		return fmt.Errorf("failed to get multiple accounts: %v", err)
	}

	return nil
}
