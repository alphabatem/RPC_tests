package methods

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// GetAccountInfo fetches the account info for a given account address
func (r *RPCTest) GetAccountInfo(accountAddress string) error {
	// Parse the account address
	pubKey, err := solana.PublicKeyFromBase58(accountAddress)
	if err != nil {
		return fmt.Errorf("invalid account address: %v", err)
	}

	// Fetch account info
	_, err = r.rpc.GetAccountInfo(
		context.Background(),
		pubKey,
	)
	if err != nil {
		return fmt.Errorf("failed to get account info: %v", err)
	}

	return nil
}
