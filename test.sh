#!/bin/bash

# Build the binary
go build -o rpc_test

# Run the test getProgramAccounts
./rpc_test getProgramAccounts --program 2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c --url http://localhost:8080 --concurrency 1 --duration 10

# Run the test getAccountInfo
./rpc_test getAccountInfo --account-file ./2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c.txt --url http://localhost:8080 --concurrency 1 --duration 10

# Run the test getMultipleAccounts
./rpc_test getMultipleAccounts --account-file ./2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c.txt --url http://localhost:8080 --concurrency 1 --duration 10
