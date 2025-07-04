# RPC Test Suite

A CLI tool for stress testing and benchmarking Solana RPC endpoints with customizable concurrency and account inputs. This tool helps developers evaluate the performance of different RPC endpoints, optimize applications, and identify potential bottlenecks.

## Features

- **Comprehensive Test Suite**: Run all RPC methods simultaneously with the `runall` command
- **Dynamic Configuration**: Generate and load test configurations with API keys
- **Smart Progress Tracking**: Real-time progress bars and detailed statistics
- **Dynamic Latency Display**: Automatic unit selection (μs, ms, s) based on performance
- **Test different Solana RPC methods** (getAccountInfo, getProgramAccounts, getMultipleAccounts)
- **Configure concurrency level** for parallel requests
- **Specify test duration**
- **Provide accounts/programs** individually or from a file
- **Seed account addresses** from programs for testing purposes
- **Comprehensive performance metrics** (requests/second, latency statistics)
- **Limit the number of accounts/programs** to process

## Installation

```bash
go build -o rpc_test
```

## Quick Start

### Comprehensive Testing with `runall`

The `runall` command provides a complete testing workflow:

```bash
# Run comprehensive test suite with API key and target RPC
./rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com

# Run with custom settings
./rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com --concurrency 10 --duration 30

# Limit accounts used for testing
./rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com --limit 50
```

**What `runall` does:**
1. **Generates test configuration** with your API key
2. **Seeds 100 accounts** from the specified program using remote RPC and gPA
3. **Runs all RPC methods** concurrently against your target RPC, using seeded accounts as needed
4. **Provides comprehensive statistics** with dynamic latency display

## Usage

### Comprehensive Testing

```bash
# Basic comprehensive test (requires --url flag)
./rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com

# Advanced comprehensive test
./rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com --concurrency 15 --duration 45 --limit 200
```

### Individual Method Testing

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

- `runall`: Execute comprehensive test suite with all methods
- `getAccountInfo`: Run tests using only the getAccountInfo RPC method
- `getMultipleAccounts`: Run tests using only the getMultipleAccounts RPC method
- `getProgramAccounts`: Run tests using only the getProgramAccounts RPC method
- `seed`: Fetch program accounts and save their addresses to a file for testing purposes
- `localRpc`: Run tests against a local RPC endpoint (e.g. lantern) running on localhost:8080 

### Global Flags (applicable to all commands)

- `-u, --url`: RPC endpoint URL (default: "https://api.mainnet-beta.solana.com")
- `-c, --concurrency`: Number of concurrent requests (default: 1)
- `-d, --duration`: Test duration in seconds (default: 10)
- `-l, --limit`: Limit the number of accounts/programs to process (0 for no limit)

### Command-specific Flags

#### runall

- `-k, --api-key`: API key for RPC endpoint (will be saved in config)
- `-c, --concurrency`: Number of concurrent requests per method (default: 5)
- `-d, --duration`: Test duration in seconds per method (default: 15)
- `-l, --limit`: Limit the number of accounts to use (0 for no limit)

**Note**: The `--url` flag is **REQUIRED** for `runall` command as it specifies the target RPC for testing.

#### getAccountInfo

- `-a, --account`: Accounts to use in tests (accepts multiple accounts, will rotate between them)
- `-f, --account-file`: File containing accounts (one per line, requests sent will rotate between them)

#### getMultipleAccounts

- `-a, --account`: Accounts to use in tests (will rotate between specified accounts in blocks of 5-15, randomly selected)
- `-f, --account-file`: File containing accounts (one per line, will rotate between them)

#### getProgramAccounts

- `-p, --program`: Program accounts to use in tests (can specify more than one)
- `-f, --program-file`: File containing program accounts (one per line)

#### seed

- `-p, --program`: Program accounts to fetch accounts from (can specify multiple programs)
- `-f, --program-file`: File containing program accounts (one per line)
- `-o, --output`: Output file to store program accounts for future tests (default: "accounts.txt")

## Dual RPC Architecture

The `runall` command uses a uses two RPCs. This lets you get program accounts from one RPC (the "remote" RPC), then use those to build tests against another (the "target"). This is mainly useful in testing Lantern, by setting Lantern as the target RPC.

### Remote RPC (Config)
- **Purpose**: Seeding account data from programs
- **Source**: Configuration file with API key
- **Use Case**: Fetching reliable account lists for testing

### Target RPC (--url flag)
- **Purpose**: Running all test methods
- **Source**: `--url` command line flag
- **Use Case**: The RPC endpoint you want to test/benchmark

This separation allows you to:
- Use a **reliable remote RPC** for getting account data
- Test any **target RPC endpoint** for performance evaluation

## Example File Format

### Account File (for getAccountInfo and getMultipleAccounts)

```
9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin
6ycRTkj1RM3L4sZcKHk8HULaFvEBaLGAiQpVpC9MPcKm
7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2
```

### Program File (for getProgramAccounts and seed)

```
2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c
FLUXubRmkEi2q6K3Y9kBPg9248ggaZVsoSFhtJHSrm1X
whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc
```

### Generated Config File (runall command)

```json
{
  "maximum_ram": 8,
  "maximum_disk": 10,
  "location": "./data/",
  "mode": "normal",
  "cache_requests": false,
  "monitoring": false,
  "monitoring_url": "",
  "log_level": "INFO",
  "rpc_url": "https://us.rpc.fluxbeam.xyz",
  "rpc_apikey": "YOUR_API_KEY_HERE",
  "programs": {
    "2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c": {
      "discriminator": 2,
      "filters": []
    }
  }
}
```

## Metrics Explanation

The test suite reports the following metrics:

### Basic Metrics
- **Total Duration**: Actual time taken to complete the test
- **Total Requests**: Number of requests processed
- **Successful Requests**: Count and percentage of successful requests
- **Failed Requests**: Count and percentage of failed requests
- **Requests per second**: Average number of requests processed per second

### Enhanced Latency Statistics (Dynamic Units)
- **Min Latency**: Minimum request latency (auto-formatted: μs, ms, or s)
- **Max Latency**: Maximum request latency (auto-formatted: μs, ms, or s)
- **Avg Latency**: Average request latency (auto-formatted: μs, ms, or s)

### Comprehensive Test Results (runall command)
- **Individual Method Results**: Detailed stats for each RPC method
- **Overall Test Summary**: Combined statistics across all methods
- **Performance Insights**: Best/worst performing methods with ratios
- **Latency Comparison**: Fastest vs slowest methods with performance ratios

## Example Output

### runall Command Output

```
🚀 Starting comprehensive RPC test suite...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 Step 1: Generating test configuration...
[██████████████████████████████] Config generated... ✅
✅ Test configuration saved to: ./config.json

📂 Step 1.5: Loading configuration with API key...
[██████████████████████████████] Config loaded... ✅
✅ Configuration loaded successfully

🌱 Step 2: Seeding accounts from program...
  🔍 Using remote RPC for seeding: https://us.rpc.fluxbeam.xyz
  🔍 Fetching accounts from program 2wT8Yq49k...
  ✅ Successfully seeded accounts
✅ Accounts seeded to: ./data/test_accounts.txt

⚡ Step 3: Running all RPC methods...
  🎯 Using target RPC for testing: https://your-target-rpc.com
  📊 Testing 3 methods with 100 accounts
  ⚙️  Concurrency: 5, Duration: 15s per method
  🔄 [1/3] Starting getAccountInfo test...
  🔄 [2/3] Starting getMultipleAccounts test...
  🔄 [3/3] Starting getProgramAccounts test...
    [████████████████░░░░] getAccountInfo: 80.0% | 12s/15s | Requests: 1250 | RPS: 104.2
    [███████████████░░░░░] getMultipleAccounts: 73.3% | 11s/15s | Requests: 1100 | RPS: 100.0
    [█████████████████░░░] getProgramAccounts: 86.7% | 13s/15s | Requests: 1300 | RPS: 100.0
    ✅ getAccountInfo completed successfully
    ✅ getMultipleAccounts completed successfully
    ✅ getProgramAccounts completed successfully

📊 Step 4: Generating comprehensive statistics...
[██████████████████████████████] Statistics calculated... ✅

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 COMPREHENSIVE TEST RESULTS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔍 INDIVIDUAL METHOD RESULTS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📈 GETACCOUNTINFO:
   Duration:         15.00 seconds
   Total Requests:    1250
   Successful:        1245 (99.60%)
   Failed:            5 (0.40%)
   Requests/second:   83.00
   Min Latency:       45.23 μs
   Max Latency:       125.67 μs
   Avg Latency:       78.45 μs

🎯 OVERALL TEST SUMMARY:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🕒 Total Duration:     45.00 seconds
🔢 Total Requests:      3650
✅ Total Successful:    3635 (99.59%)
❌ Total Failed:        15 (0.41%)
⚡ Overall RPS:         81.11
📊 Methods Tested:      3

💡 PERFORMANCE INSIGHTS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏆 Best Performing:    getAccountInfo (83.00 RPS)
🐌 Worst Performing:   getProgramAccounts (80.00 RPS)
📊 Performance Ratio:  96.4% (worst/best)

⏱️  LATENCY COMPARISON:
⚡ Fastest Method:     getAccountInfo (78.45 μs avg)
🐌 Slowest Method:     getProgramAccounts (245.12 μs avg)
📊 Latency Ratio:      3.1x (slowest/fastest)

✅ Comprehensive test suite completed successfully!
```

## Use Cases

- **Compare performance** between different Solana RPC providers
- **Optimize application settings** when using Lantern
- **Identify performance bottlenecks** when using Lantern
- **Test the stability** of Lantern or RPC endpoints under load
