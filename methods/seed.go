package methods

import (
	"context"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
)

// SeedProgramAccounts fetches program accounts and saves their addresses to the specified output file
func (r *RPCTest) SeedProgramAccounts(programAddress string, outputFile string, limit int) error {
	// Parse the program address
	pubKey, err := solana.PublicKeyFromBase58(programAddress)
	if err != nil {
		return fmt.Errorf("invalid program address: %v", err)
	}

	// Fetch program accounts
	accounts, err := r.rpc.GetProgramAccounts(
		context.Background(),
		pubKey,
	)
	if err != nil {
		return fmt.Errorf("failed to get program accounts: %v", err)
	}

	// Create the output file if it doesn't exist
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Apply limit if specified
	totalAccounts := len(accounts)
	if limit > 0 && limit < totalAccounts {
		accounts = accounts[:limit]
		fmt.Printf("Limiting to %d accounts out of %d found for program %s\n", limit, totalAccounts, programAddress)
	} else {
		fmt.Printf("Found %d accounts for program %s\n", totalAccounts, programAddress)
	}

	fmt.Printf("Saving account addresses to %s\n", outputFile)

	// Save each account address to the file
	for i, account := range accounts {
		// Write account address to the file
		if _, err := file.WriteString(account.Pubkey.String() + "\n"); err != nil {
			return fmt.Errorf("failed to write to output file: %v", err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("Processed %d/%d accounts\n", i+1, len(accounts))
		}
	}

	fmt.Printf("Total accounts saved: %d\n", len(accounts))
	fmt.Printf("Account addresses saved to: %s\n", outputFile)
	fmt.Printf("Use this file with other commands: --account-file %s\n", outputFile)

	return nil
}
