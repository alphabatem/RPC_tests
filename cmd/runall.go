package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"rpc_test/methods"

	"github.com/spf13/cobra"
)

// TestConfig represents the configuration for the test
type TestConfig struct {
	MaximumRAM    int                    `json:"maximum_ram"`
	MaximumDisk   int                    `json:"maximum_disk"`
	Location      string                 `json:"location"`
	Mode          string                 `json:"mode"`
	CacheRequests bool                   `json:"cache_requests"`
	Monitoring    bool                   `json:"monitoring"`
	MonitoringURL string                 `json:"monitoring_url"`
	LogLevel      string                 `json:"log_level"`
	RPCURL        string                 `json:"rpc_url"`
	RPCAPIKey     string                 `json:"rpc_apikey"`
	Programs      map[string]ProgramInfo `json:"programs"`
}

// ProgramInfo represents program-specific configuration
type ProgramInfo struct {
	Discriminator int      `json:"discriminator"`
	Filters       []string `json:"filters"`
}

// TestResult represents the result of a single method test
type TestResult struct {
	MethodName     string
	Duration       time.Duration
	TotalRequests  int64
	SuccessCount   int64
	FailureCount   int64
	RequestsPerSec float64
	SuccessRate    float64
	MinLatency     time.Duration
	MaxLatency     time.Duration
	AvgLatency     time.Duration
}

// OverallResult represents the overall test results
type OverallResult struct {
	TotalDuration      time.Duration
	TotalRequests      int64
	TotalSuccess       int64
	TotalFailure       int64
	OverallRPS         float64
	OverallSuccessRate float64
	MethodResults      []TestResult
}

// Default configuration as specified
var defaultConfig = TestConfig{
	MaximumRAM:    8,
	MaximumDisk:   10,
	Location:      "./data/",
	Mode:          "normal",
	CacheRequests: false,
	Monitoring:    false,
	MonitoringURL: "",
	LogLevel:      "INFO",
	RPCURL:        "https://us.rpc.fluxbeam.xyz",
	RPCAPIKey:     "YOUR_API_KEY_HERE",
	Programs: map[string]ProgramInfo{
		"2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c": {
			Discriminator: 2,
			Filters:       []string{},
		},
	},
}

// runallCmd represents the runall command
var runallCmd = &cobra.Command{
	Use:   "runall",
	Short: "Run all methods with comprehensive testing",
	Long: `Execute a comprehensive test suite that includes:
1. Seeding data from default account configuration
2. Seeding 100 accounts from the specified program
3. Running all available RPC methods
4. Providing detailed statistics for each method and overall results

Example:
  rpc_test runall --concurrency 10 --duration 30`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting comprehensive RPC test suite...")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Step 1: Generate and save test configuration
		fmt.Println("📋 Step 1: Generating test configuration...")
		showProgress("Generating config", 100)
		configFile := "./config.json"
		if err := generateTestConfig(configFile); err != nil {
			log.Fatalf("Failed to generate test config: %v", err)
		}
		showProgressComplete("Config generated")
		fmt.Printf("✅ Test configuration saved to: %s\n", configFile)

		// Step 1.5: Load the generated config with API key
		fmt.Println("\n📂 Step 1.5: Loading configuration with API key...")
		showProgress("Loading config", 100)
		config, err := loadTestConfig(configFile)
		if err != nil {
			log.Fatalf("Failed to load test config: %v", err)
		}
		showProgressComplete("Config loaded")
		fmt.Printf("✅ Configuration loaded successfully\n")

		// Step 2: Seed accounts from the program
		fmt.Println("\n🌱 Step 2: Seeding accounts from program...")
		accountsFile := "./data/test_accounts.txt"
		if err := seedAccountsFromProgram(accountsFile, config); err != nil {
			log.Fatalf("Failed to seed accounts: %v", err)
		}
		fmt.Printf("✅ Accounts seeded to: %s\n", accountsFile)

		// Step 3: Run all methods
		fmt.Println("\n⚡ Step 3: Running all RPC methods...")
		results, err := runAllMethods(accountsFile)
		if err != nil {
			log.Fatalf("Failed to run methods: %v", err)
		}

		// Step 4: Generate and display statistics
		fmt.Println("\n📊 Step 4: Generating comprehensive statistics...")
		showProgress("Calculating statistics", 100)
		overallResult := calculateOverallResults(results)
		showProgressComplete("Statistics calculated")
		displayResults(results, overallResult)
	},
}

// showProgress displays a progress bar with the given message and percentage
func showProgress(message string, percentage int) {
	const barWidth = 30
	progress := int(float64(percentage) * float64(barWidth) / 100)
	progressBar := strings.Repeat("█", progress) + strings.Repeat("░", barWidth-progress)
	fmt.Printf("\r[%s] %s... %d%%", progressBar, message, percentage)
}

// showProgressComplete displays a completed progress bar
func showProgressComplete(message string) {
	const barWidth = 30
	progressBar := strings.Repeat("█", barWidth)
	fmt.Printf("\r[%s] %s... ✅\n", progressBar, message)
}

// generateTestConfig creates and saves the test configuration
func generateTestConfig(configFile string) error {
	// Create data directory if it doesn't exist
	dataDir := filepath.Dir(configFile)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Create a copy of the default config to modify
	config := defaultConfig

	// Use provided API key if available
	if apiKey != "" {
		config.RPCAPIKey = apiKey
		fmt.Printf("✅ Using provided API key: %s...\n", apiKey[:8]+"***")
	} else {
		fmt.Println("⚠️  WARNING: No API key provided!")
		fmt.Println("   Please edit the generated config file to set your API key:")
		fmt.Printf("   %s\n", configFile)
		fmt.Println("   Or use the --api-key flag to provide it directly.")
	}

	// Marshal configuration to JSON
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write configuration to file
	if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// loadTestConfig loads the test configuration from file
func loadTestConfig(configFile string) (TestConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return TestConfig{}, fmt.Errorf("failed to read config file: %v", err)
	}

	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return TestConfig{}, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// seedAccountsFromProgram seeds 100 accounts from the default program
func seedAccountsFromProgram(accountsFile string, config TestConfig) error {
	// Create data directory if it doesn't exist
	dataDir := filepath.Dir(accountsFile)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Get the program ID from default config
	var programID string
	for program := range config.Programs {
		programID = program
		break
	}

	if programID == "" {
		return fmt.Errorf("no program found in default configuration")
	}

	// Use the config RPC URL for seeding (remote RPC)
	seedRPCURL := config.RPCURL
	if config.RPCAPIKey != "" && config.RPCAPIKey != "YOUR_API_KEY_HERE" {
		seedRPCURL = fmt.Sprintf("%s?key=%s", config.RPCURL, config.RPCAPIKey)
	}

	fmt.Printf("  🔍 Using remote RPC for seeding: %s\n", config.RPCURL)
	fmt.Printf("  🔍 Fetching accounts from program %s...\n", programID[:8]+"...")

	// Create RPC client for seeding (using config RPC URL)
	rpcTest := methods.NewRPCTest(seedRPCURL)

	// Seed program accounts with limit of 100
	err := rpcTest.SeedProgramAccounts(programID, accountsFile, 100)
	if err != nil {
		return err
	}

	// Show completion
	fmt.Printf("  ✅ Successfully seeded accounts\n")
	return nil
}

// runAllMethods runs all available RPC methods and returns results
func runAllMethods(accountsFile string) ([]TestResult, error) {
	// Check if --url flag is provided for target RPC
	if rpcURL == "" || rpcURL == "https://api.mainnet-beta.solana.com" {
		log.Fatalf("❌ ERROR: --url flag is required for target RPC testing!")
		fmt.Println("   Please provide the target RPC endpoint using --url flag.")
		fmt.Println("   Example: --url https://your-target-rpc.com")
		fmt.Println("   This is the RPC endpoint you want to test/benchmark.")
	}

	fmt.Printf("  🎯 Using target RPC for testing: %s\n", rpcURL)

	// Define all available methods
	methods := []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"}

	// Load accounts from file
	data, err := os.ReadFile(accountsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var accounts []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			accounts = append(accounts, line)
		}
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found in file")
	}

	// Apply limit if specified
	if limit > 0 && limit < len(accounts) {
		accounts = accounts[:limit]
	}

	fmt.Printf("  📊 Testing %d methods with %d accounts\n", len(methods), len(accounts))
	fmt.Printf("  ⚙️  Concurrency: %d, Duration: %ds per method\n", concurrency, duration)

	var results []TestResult
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Run each method concurrently
	for i, methodName := range methods {
		wg.Add(1)
		go func(method string, methodIndex int) {
			defer wg.Done()

			result := runSingleMethod(method, accounts, methodIndex+1, len(methods))

			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(methodName, i)
	}

	wg.Wait()

	return results, nil
}

// runSingleMethod runs a single method test and returns the result
func runSingleMethod(methodName string, accounts []string, methodIndex, totalMethods int) TestResult {
	fmt.Printf("  🔄 [%d/%d] Starting %s test...\n", methodIndex, totalMethods, methodName)

	// Create RPC client with target RPC URL (from --url flag)
	rpcTest := methods.NewRPCTest(rpcURL)

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	var wg sync.WaitGroup
	var successCount, failureCount int64
	var mutex sync.Mutex

	// Collect statistics
	var totalLatency time.Duration
	var minLatency time.Duration = time.Hour
	var maxLatency time.Duration

	// Create channels for workers
	stop := make(chan struct{})

	// Progress reporting
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Now().After(endTime) {
					return
				}

				mutex.Lock()
				elapsed := time.Since(startTime)
				currentTotal := successCount + failureCount
				currentRPS := float64(currentTotal) / elapsed.Seconds()
				percentComplete := (elapsed.Seconds() / float64(duration)) * 100

				// Create a simple progress bar
				const barWidth = 20
				progress := int(percentComplete * float64(barWidth) / 100)
				progressBar := strings.Repeat("█", progress) + strings.Repeat("░", barWidth-progress)

				fmt.Printf("\r    [%s] %s: %.1f%% | %ds/%ds | Requests: %d | RPS: %.1f",
					progressBar, methodName, percentComplete, int(elapsed.Seconds()), duration, currentTotal, currentRPS)
				mutex.Unlock()
			case <-stop:
				return
			}
		}
	}()

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-stop:
					return
				default:
					// Check if test duration has elapsed
					if time.Now().After(endTime) {
						return
					}

					// Execute the specified method
					startReq := time.Now()
					err := Method(methodName, rpcTest, accounts[workerID%len(accounts)])
					reqDuration := time.Since(startReq)

					mutex.Lock()
					if err != nil {
						failureCount++
					} else {
						successCount++
						totalLatency += reqDuration
						if reqDuration < minLatency {
							minLatency = reqDuration
						}
						if reqDuration > maxLatency {
							maxLatency = reqDuration
						}
					}
					mutex.Unlock()
				}
			}
		}(i)
	}

	// Wait for the test duration
	time.Sleep(time.Duration(duration) * time.Second)
	close(stop)

	// Wait for all workers to finish
	wg.Wait()

	// Clear the progress line and show completion
	fmt.Printf("\r    ✅ %s completed successfully\n", methodName)

	// Calculate results
	totalDuration := time.Since(startTime)
	totalRequests := successCount + failureCount
	requestsPerSecond := float64(totalRequests) / totalDuration.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100

	var avgLatency time.Duration
	if successCount > 0 {
		avgLatency = totalLatency / time.Duration(successCount)
	}

	return TestResult{
		MethodName:     methodName,
		Duration:       totalDuration,
		TotalRequests:  totalRequests,
		SuccessCount:   successCount,
		FailureCount:   failureCount,
		RequestsPerSec: requestsPerSecond,
		SuccessRate:    successRate,
		MinLatency:     minLatency,
		MaxLatency:     maxLatency,
		AvgLatency:     avgLatency,
	}
}

// calculateOverallResults calculates overall statistics
func calculateOverallResults(methodResults []TestResult) OverallResult {
	var totalDuration time.Duration
	var totalRequests, totalSuccess, totalFailure int64

	for _, result := range methodResults {
		totalDuration += result.Duration
		totalRequests += result.TotalRequests
		totalSuccess += result.SuccessCount
		totalFailure += result.FailureCount
	}

	overallRPS := float64(totalRequests) / totalDuration.Seconds()
	overallSuccessRate := float64(totalSuccess) / float64(totalRequests) * 100

	return OverallResult{
		TotalDuration:      totalDuration,
		TotalRequests:      totalRequests,
		TotalSuccess:       totalSuccess,
		TotalFailure:       totalFailure,
		OverallRPS:         overallRPS,
		OverallSuccessRate: overallSuccessRate,
		MethodResults:      methodResults,
	}
}

// formatLatency formats latency in the most appropriate unit
func formatLatency(duration time.Duration) string {
	if duration < time.Millisecond {
		return fmt.Sprintf("%.2f μs", float64(duration.Microseconds()))
	} else if duration < time.Second {
		return fmt.Sprintf("%.2f ms", float64(duration.Milliseconds()))
	} else {
		return fmt.Sprintf("%.2f s", duration.Seconds())
	}
}

// displayResults displays comprehensive test results
func displayResults(methodResults []TestResult, overall OverallResult) {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 COMPREHENSIVE TEST RESULTS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Display individual method results
	fmt.Println("\n🔍 INDIVIDUAL METHOD RESULTS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, result := range methodResults {
		fmt.Printf("\n📈 %s:\n", strings.ToUpper(result.MethodName))
		fmt.Printf("   Duration:         %.2f seconds\n", result.Duration.Seconds())
		fmt.Printf("   Total Requests:    %d\n", result.TotalRequests)
		fmt.Printf("   Successful:        %d (%.2f%%)\n", result.SuccessCount, result.SuccessRate)
		fmt.Printf("   Failed:            %d (%.2f%%)\n", result.FailureCount, 100-result.SuccessRate)
		fmt.Printf("   Requests/second:   %.2f\n", result.RequestsPerSec)
		if result.SuccessCount > 0 {
			fmt.Printf("   Min Latency:       %s\n", formatLatency(result.MinLatency))
			fmt.Printf("   Max Latency:       %s\n", formatLatency(result.MaxLatency))
			fmt.Printf("   Avg Latency:       %s\n", formatLatency(result.AvgLatency))
		}
	}

	// Display overall results
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 OVERALL TEST SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🕒 Total Duration:     %.2f seconds\n", overall.TotalDuration.Seconds())
	fmt.Printf("🔢 Total Requests:      %d\n", overall.TotalRequests)
	fmt.Printf("✅ Total Successful:    %d (%.2f%%)\n", overall.TotalSuccess, overall.OverallSuccessRate)
	fmt.Printf("❌ Total Failed:        %d (%.2f%%)\n", overall.TotalFailure, 100-overall.OverallSuccessRate)
	fmt.Printf("⚡ Overall RPS:         %.2f\n", overall.OverallRPS)
	fmt.Printf("📊 Methods Tested:      %d\n", len(methodResults))

	// Performance insights
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 PERFORMANCE INSIGHTS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Find best and worst performing methods
	var bestMethod, worstMethod TestResult
	var bestRPS, worstRPS float64

	for i, result := range methodResults {
		if i == 0 || result.RequestsPerSec > bestRPS {
			bestMethod = result
			bestRPS = result.RequestsPerSec
		}
		if i == 0 || result.RequestsPerSec < worstRPS {
			worstMethod = result
			worstRPS = result.RequestsPerSec
		}
	}

	fmt.Printf("🏆 Best Performing:    %s (%.2f RPS)\n", bestMethod.MethodName, bestRPS)
	fmt.Printf("🐌 Worst Performing:   %s (%.2f RPS)\n", worstMethod.MethodName, worstRPS)

	if bestRPS > 0 {
		performanceRatio := worstRPS / bestRPS * 100
		fmt.Printf("📊 Performance Ratio:  %.1f%% (worst/best)\n", performanceRatio)
	}

	// Add latency comparison
	if len(methodResults) > 0 {
		fmt.Println("\n⏱️  LATENCY COMPARISON:")
		var fastestMethod, slowestMethod TestResult
		var fastestLatency, slowestLatency time.Duration

		for i, result := range methodResults {
			if result.SuccessCount > 0 {
				if i == 0 || result.AvgLatency < fastestLatency {
					fastestMethod = result
					fastestLatency = result.AvgLatency
				}
				if i == 0 || result.AvgLatency > slowestLatency {
					slowestMethod = result
					slowestLatency = result.AvgLatency
				}
			}
		}

		if fastestLatency > 0 {
			fmt.Printf("⚡ Fastest Method:     %s (%s avg)\n", fastestMethod.MethodName, formatLatency(fastestLatency))
			fmt.Printf("🐌 Slowest Method:     %s (%s avg)\n", slowestMethod.MethodName, formatLatency(slowestLatency))

			if fastestLatency > 0 {
				latencyRatio := float64(slowestLatency) / float64(fastestLatency)
				fmt.Printf("📊 Latency Ratio:      %.1fx (slowest/fastest)\n", latencyRatio)
			}
		}
	}

	fmt.Println("\n✅ Comprehensive test suite completed successfully!")
}

func init() {
	RootCmd.AddCommand(runallCmd)

	// Add runall-specific flags
	runallCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 5, "Number of concurrent requests per method")
	runallCmd.Flags().IntVarP(&duration, "duration", "d", 15, "Test duration in seconds per method")
	runallCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit the number of accounts to use (0 for no limit)")
	runallCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for RPC endpoint (will be saved in config)")
}
