package methods

import (
	"github.com/gagliardetto/solana-go/rpc"
)

type RPCTest struct {
	rpc    *rpc.Client
	rpcUrl string
}

func NewRPCTest(rpcUrl string) *RPCTest {
	return &RPCTest{rpc: rpc.New(rpcUrl), rpcUrl: rpcUrl}
}
