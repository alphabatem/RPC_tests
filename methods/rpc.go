package methods

import (
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
)

type RPCTest struct {
	rpc    *rpc.Client
	rpcUrl string
}

func NewRPCTest(rpcUrl string, apiKey string) *RPCTest {
	url := fmt.Sprintf("%s?key=%s", rpcUrl, apiKey)
	return &RPCTest{rpc: rpc.New(url), rpcUrl: url}
}
