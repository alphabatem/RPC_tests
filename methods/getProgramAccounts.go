package methods

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// GetProgramAccounts fetches accounts owned by the program
func (r *RPCTest) GetProgramAccounts(programAddress string) error {
	// Parse the program address
	pubKey, err := solana.PublicKeyFromBase58(programAddress)
	if err != nil {
		return fmt.Errorf("invalid program address: %v", err)
	}

	// Fetch program accounts
	_, err = r.rpc.GetProgramAccounts(
		context.Background(),
		pubKey,
	)
	if err != nil {
		return fmt.Errorf("failed to get program accounts: %v", err)
	}

	return nil
}
