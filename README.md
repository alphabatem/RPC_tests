# RPC Test Suite

A CLI tool for stress testing and benchmarking Solana RPC endpoints with customizable concurrency and account inputs. This tool helps developers evaluate the performance of different RPC endpoints, optimize applications, and identify potential bottlenecks.

## Features

- Test different Solana RPC methods (getAccountInfo, getProgramAccounts, getMultipleAccounts)
- Configure concurrency level for parallel requests
- Specify test duration
- Provide accounts/programs individually or from a file
- Seed account addresses from programs for testing purposes
- Test against both public and local RPC endpoints
- Comprehensive performance metrics (requests/second, latency statistics)
- Limit the number of accounts/programs to process

## Installation

```bash
go build -o rpc_test
```

## Usage

```bash
# Basic usage for getAccountInfo
./rpc_test getAccountInfo --account <ACCOUNT_ADDRESS> --concurrency 10 --duration 30

# Using an account file for getAccountInfo
./rpc_test getAccountInfo --account-file accounts.txt --concurrency 10 --duration 30

# Basic usage for getMultipleAccounts
./rpc_test getMultipleAccounts --account <ACCOUNT_ADDRESS> --concurrency 10 --duration 30

# Using comma-separated accounts for getMultipleAccounts
./rpc_test getMultipleAccounts --account "addr1,addr2,addr3" --concurrency 10 --duration 30

# Basic usage for getProgramAccounts
./rpc_test getProgramAccounts --program <PROGRAM_ADDRESS> --concurrency 50 --duration 60

# Using a program file for getProgramAccounts
./rpc_test getProgramAccounts --program-file programs.txt --concurrency 50 --duration 60

# Seed account data from a program for testing
./rpc_test seed --program <PROGRAM_ADDRESS> --output accounts.txt

# Seed account data with a limit of 1000 accounts
./rpc_test seed --program <PROGRAM_ADDRESS> --output accounts.txt --limit 1000

# Test against a local RPC endpoint
./rpc_test localRpc getAccountInfo --account <ACCOUNT_ADDRESS> --concurrency 5

# Test with limited number of accounts
./rpc_test getAccountInfo --account-file accounts.txt --limit 100 --concurrency 10
```

### Available Commands

- `getAccountInfo`: Run tests against the getAccountInfo RPC method
- `getMultipleAccounts`: Run tests against the getMultipleAccounts RPC method
- `getProgramAccounts`: Run tests against the getProgramAccounts RPC method
- `seed`: Fetch program accounts and save their addresses to a file for testing purposes
- `localRpc`: Run tests against a local RPC endpoint running on localhost:8080

### Global Flags (applicable to all commands)

- `-u, --url`: RPC endpoint URL (default: "https://api.mainnet-beta.solana.com")
- `-c, --concurrency`: Number of concurrent requests (default: 1)
- `-d, --duration`: Test duration in seconds (default: 10)
- `-l, --limit`: Limit the number of accounts/programs to process (0 for no limit)

### Command-specific Flags

#### getAccountInfo

- `-a, --account`: Account addresses to use in tests (can be specified multiple times)
- `-f, --account-file`: File containing account addresses (one per line)

#### getMultipleAccounts

- `-a, --account`: Account addresses to use in tests (can be specified multiple times or comma-separated)
- `-f, --account-file`: File containing account addresses (one per line)

#### getProgramAccounts

- `-p, --program`: Program addresses to use in tests (can be specified multiple times)
- `-f, --program-file`: File containing program addresses (one per line)

#### seed

- `-p, --program`: Program addresses to fetch accounts for (can be specified multiple times)
- `-f, --program-file`: File containing program addresses (one per line)
- `-o, --output`: Output file to store account addresses (default: "accounts.txt")

## Example File Format

### Account File (for getAccountInfo and getMultipleAccounts)

```
9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin
6ycRTkj1RM3L4sZcKHk8HULaFvEBaLGAiQpVpC9MPcKm
7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2
```

### Program File (for getProgramAccounts and seed)

```
TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA
ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL
ComputeBudget111111111111111111111111111111
```

## Metrics Explanation

The test suite reports the following metrics:

- **Total Duration**: Actual time taken to complete the test
- **Total Requests**: Number of requests processed
- **Successful Requests**: Count and percentage of successful requests
- **Failed Requests**: Count and percentage of failed requests
- **Requests per second**: Average number of requests processed per second
- **Latency Statistics**: 
  - Min: Minimum request latency in milliseconds
  - Max: Maximum request latency in milliseconds
  - Avg: Average request latency in milliseconds

## Use Cases

- Compare performance between different Solana RPC providers
- Optimize application settings for RPC requests
- Identify performance bottlenecks in RPC interactions
- Test the stability of RPC endpoints under load
- Benchmark local Solana validator nodes 