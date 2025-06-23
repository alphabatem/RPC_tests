# RPC Test Suite

A comprehensive CLI tool for stress testing and benchmarking Solana RPC endpoints with customizable concurrency and account inputs. This tool helps developers evaluate the performance of different RPC endpoints, optimize applications, and identify potential bottlenecks in Solana blockchain interactions.

## 🚀 Features

- **Comprehensive Test Suite**: Run all RPC methods simultaneously with the `runall` command
- **Dual RPC Architecture**: Use remote RPC for seeding accounts and target RPC for testing
- **Dynamic Configuration**: Generate and load test configurations with API keys
- **Smart Progress Tracking**: Real-time progress bars and detailed statistics
- **Dynamic Latency Display**: Automatic unit selection (μs, ms, s) based on performance
- **Test different Solana RPC methods** (getAccountInfo, getProgramAccounts, getMultipleAccounts)
- **Configure concurrency level** for parallel requests
- **Specify test duration**
- **Provide accounts/programs** individually or from a file
- **Seed account addresses** from programs for testing purposes
- **Test against both public and local RPC endpoints**
- **Comprehensive performance metrics** (requests/second, latency statistics)
- **Limit the number of accounts/programs** to process
- **Real-time progress monitoring** with visual progress bars
- **Configurable test parameters** for different testing scenarios

## 📋 Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Architecture Overview](#architecture-overview)
- [Usage](#usage)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [Contributing](#contributing)

## 🛠️ Installation

### Prerequisites

- Go 1.24.2 or higher
- Git

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd rpc_test

# Build the application
go build -o rpc_test

# Make it executable (Linux/macOS)
chmod +x rpc_test
```

### Verify Installation

```bash
# Check if the binary was created successfully
./rpc_test --help
```

## 🚀 Quick Start

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
2. **Seeds 100 accounts** from the specified program using remote RPC
3. **Runs all RPC methods** concurrently against your target RPC
4. **Provides comprehensive statistics** with dynamic latency display

## 📁 Project Structure

```
rpc_test/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── README.md              # This file
├── config-template.json   # Template configuration file
├── config.json            # Generated configuration (gitignored)
├── .gitignore             # Git ignore rules
├── rpc_test               # Compiled binary (gitignored)
├── cmd/                   # Command implementations
│   ├── root.go           # Root command and global flags
│   ├── common.go         # Shared utilities and variables
│   ├── runall.go         # Comprehensive test suite command
│   ├── getAccountInfo.go # getAccountInfo RPC testing
│   ├── getMultipleAccounts.go # getMultipleAccounts RPC testing
│   ├── getProgramAccounts.go # getProgramAccounts RPC testing
│   └── seed.go           # Account seeding functionality
├── methods/               # RPC method implementations
│   ├── rpc.go            # Base RPC client wrapper
│   ├── getAccountInfo.go # getAccountInfo implementation
│   ├── getMultipleAccounts.go # getMultipleAccounts implementation
│   ├── getProgramAccounts.go # getProgramAccounts implementation
│   └── seed.go           # Account seeding logic
├── data/                  # Test data and generated files
│   └── test_accounts.txt # Generated test accounts
└── *.txt                 # Example account/program files
```

## 🏗️ Architecture Overview

### Core Components

1. **Command Layer** (`cmd/`): Implements CLI commands using Cobra framework
2. **Method Layer** (`methods/`): Contains RPC method implementations using solana-go
3. **Configuration Management**: Dynamic config generation and loading
4. **Progress Tracking**: Real-time progress monitoring with visual feedback
5. **Statistics Engine**: Comprehensive metrics calculation and reporting

### Dual RPC Architecture

The application uses a sophisticated dual RPC architecture for optimal performance:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Remote RPC    │    │   Target RPC    │    │   Local Files   │
│   (Config)      │    │   (--url flag)  │    │   (data/)       │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • Seeding       │    │ • Testing       │    │ • Account lists │
│ • Data fetch    │    │ • Benchmarking  │    │ • Config files  │
│ • Reliable      │    │ • Performance   │    │ • Results       │
│ • API key auth  │    │ • Load testing  │    │ • Logs          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Data Flow

1. **Configuration Phase**: Generate/load config with API keys
2. **Seeding Phase**: Fetch account data from reliable remote RPC
3. **Testing Phase**: Run tests against target RPC endpoint
4. **Analysis Phase**: Calculate and display comprehensive statistics

## 📖 Usage

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

- `runall`: **NEW** - Execute comprehensive test suite with all methods
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

#### runall

- `-k, --api-key`: API key for RPC endpoint (will be saved in config)
- `-c, --concurrency`: Number of concurrent requests per method (default: 5)
- `-d, --duration`: Test duration in seconds per method (default: 15)
- `-l, --limit`: Limit the number of accounts to use (0 for no limit)

**Note**: The `--url` flag is **REQUIRED** for `runall` command as it specifies the target RPC for testing.

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

## ⚙️ Configuration

### Configuration File Structure

The application generates a `config.json` file with the following structure:

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

### Environment Variables

The application respects the following environment variables:

- `RPC_API_KEY`: Default API key for RPC endpoints
- `RPC_URL`: Default RPC endpoint URL
- `CONCURRENCY`: Default concurrency level
- `DURATION`: Default test duration

### Configuration Management

1. **Auto-generation**: The `runall` command automatically generates configuration
2. **API Key Storage**: API keys are securely stored in the config file
3. **Template-based**: Uses `config-template.json` as a base template
4. **Dynamic Loading**: Configuration is loaded at runtime

## 📊 API Reference

### RPC Methods Supported

#### getAccountInfo
- **Purpose**: Fetch account information for specific addresses
- **Use Case**: Testing account data retrieval performance
- **Parameters**: Account addresses (single or multiple)

#### getMultipleAccounts
- **Purpose**: Fetch information for multiple accounts in a single request
- **Use Case**: Testing batch account data retrieval
- **Parameters**: Multiple account addresses

#### getProgramAccounts
- **Purpose**: Fetch all accounts owned by a specific program
- **Use Case**: Testing program account enumeration
- **Parameters**: Program addresses

### Performance Metrics

The test suite reports comprehensive metrics:

#### Basic Metrics
- **Total Duration**: Actual time taken to complete the test
- **Total Requests**: Number of requests processed
- **Successful Requests**: Count and percentage of successful requests
- **Failed Requests**: Count and percentage of failed requests
- **Requests per second**: Average number of requests processed per second

#### Enhanced Latency Statistics (Dynamic Units)
- **Min Latency**: Minimum request latency (auto-formatted: μs, ms, or s)
- **Max Latency**: Maximum request latency (auto-formatted: μs, ms, or s)
- **Avg Latency**: Average request latency (auto-formatted: μs, ms, or s)

#### Comprehensive Test Results (runall command)
- **Individual Method Results**: Detailed stats for each RPC method
- **Overall Test Summary**: Combined statistics across all methods
- **Performance Insights**: Best/worst performing methods with ratios
- **Latency Comparison**: Fastest vs slowest methods with performance ratios

## 📝 Examples

### Example File Format

#### Account File (for getAccountInfo and getMultipleAccounts)

```
9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin
6ycRTkj1RM3L4sZcKHk8HULaFvEBaLGAiQpVpC9MPcKm
7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2
```

#### Program File (for getProgramAccounts and seed)

```
TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA
ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL
ComputeBudget111111111111111111111111111111
```

### Example Output

#### runall Command Output

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

## 🔧 Troubleshooting

### Common Issues

#### 1. API Key Authentication Errors
```
Error: authentication failed
```
**Solution**: Ensure your API key is valid and has the necessary permissions.

#### 2. RPC Endpoint Connection Issues
```
Error: connection refused
```
**Solution**: Verify the RPC endpoint URL and check network connectivity.

#### 3. Rate Limiting
```
Error: rate limit exceeded
```
**Solution**: Reduce concurrency or increase delays between requests.

#### 4. Insufficient Account Data
```
Error: no accounts found
```
**Solution**: Check if the program address is valid and contains accounts.

### Debug Mode

Enable verbose logging for debugging:

```bash
# Set log level to DEBUG
export LOG_LEVEL=DEBUG
./rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com
```

### Performance Optimization

1. **Concurrency Tuning**: Start with low concurrency and gradually increase
2. **Duration Adjustment**: Use longer durations for more accurate metrics
3. **Account Limits**: Limit accounts for faster testing iterations
4. **Network Optimization**: Use RPC endpoints closer to your location

## 🛠️ Development

### Prerequisites

- Go 1.24.2+
- Git
- Make (optional, for build scripts)

### Development Setup

```bash
# Clone the repository
git clone <repository-url>
cd rpc_test

# Install dependencies
go mod download

# Run tests
go test ./...

# Build for development
go build -o rpc_test

# Run with development flags
./rpc_test --help
```

### Project Dependencies

Key dependencies used in this project:

- **github.com/spf13/cobra**: CLI framework
- **github.com/gagliardetto/solana-go**: Solana blockchain interaction
- **go.uber.org/zap**: Structured logging
- **github.com/fatih/color**: Colored terminal output

### Code Structure

#### Command Layer (`cmd/`)
- **root.go**: Main command setup and global flags
- **common.go**: Shared utilities and variables
- **runall.go**: Comprehensive test suite implementation
- **getAccountInfo.go**: getAccountInfo command
- **getMultipleAccounts.go**: getMultipleAccounts command
- **getProgramAccounts.go**: getProgramAccounts command
- **seed.go**: Account seeding command

#### Method Layer (`methods/`)
- **rpc.go**: Base RPC client wrapper
- **getAccountInfo.go**: getAccountInfo RPC implementation
- **getMultipleAccounts.go**: getMultipleAccounts RPC implementation
- **getProgramAccounts.go**: getProgramAccounts RPC implementation
- **seed.go**: Account seeding logic

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./methods -v

# Run benchmarks
go test -bench=. ./...
```

### Building

```bash
# Build for current platform
go build -o rpc_test

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o rpc_test_linux
GOOS=darwin GOARCH=amd64 go build -o rpc_test_macos
GOOS=windows GOARCH=amd64 go build -o rpc_test_windows.exe
```

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `go test ./...`
6. Commit your changes: `git commit -am 'Add new feature'`
7. Push to the branch: `git push origin feature/your-feature`
8. Submit a pull request

### Code Style

- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for complex logic
- Ensure all tests pass
- Update documentation as needed

### Testing Guidelines

- Write unit tests for new functionality
- Ensure existing tests continue to pass
- Add integration tests for new commands
- Test with different RPC endpoints

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Solana Labs for the Solana blockchain
- Gagliardetto for the excellent solana-go library
- The Solana community for feedback and contributions

## 📞 Support

For support and questions:

- Create an issue on GitHub
- Check the troubleshooting section
- Review the examples and documentation

---

**Note**: This tool is designed for testing and benchmarking purposes. Please use responsibly and respect RPC endpoint rate limits and terms of service. 