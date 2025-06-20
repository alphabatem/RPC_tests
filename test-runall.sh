#!/bin/bash

# Build the binary
go build -o rpc_test

# Run the test getProgramAccounts
./rpc_test runall --api-key <YOUR_API_KEY> --concurrency 1 --duration 10 --limit 100 -u http://localhost:8080