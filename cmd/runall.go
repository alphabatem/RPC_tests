package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	RemoteRPCURL string   `json:"rpc_url"`
	RPCAPIKey    string   `json:"rpc_apikey"`
	Programs     []string `json:"programs"`
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
	RemoteRPCURL: "https://us.rpc.fluxbeam.xyz",
	RPCAPIKey:    "YOUR_API_KEY_HERE",
	Programs:     []string{"2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"},
}

// ProgressManager manages progress display for all methods
type ProgressManager struct {
	methods      map[string]*MethodProgress
	mutex        sync.RWMutex
	stopChan     chan struct{}
	firstDisplay bool
}

// MethodProgress tracks progress for a single method
type MethodProgress struct {
	Name            string
	StartTime       time.Time
	EndTime         time.Time
	SuccessCount    int64
	FailureCount    int64
	TotalRequests   int64
	RequestsPerSec  float64
	PercentComplete float64
}

// NewProgressManager creates a new progress manager
func NewProgressManager() *ProgressManager {
	return &ProgressManager{
		methods:      make(map[string]*MethodProgress),
		stopChan:     make(chan struct{}),
		firstDisplay: true,
	}
}

// RegisterMethod registers a method for progress tracking
func (pm *ProgressManager) RegisterMethod(methodName string, duration int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	pm.methods[methodName] = &MethodProgress{
		Name:      methodName,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// UpdateProgress updates progress for a specific method
func (pm *ProgressManager) UpdateProgress(methodName string, successCount, failureCount int64) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if method, exists := pm.methods[methodName]; exists {
		method.SuccessCount = successCount
		method.FailureCount = failureCount
		method.TotalRequests = successCount + failureCount

		elapsed := time.Since(method.StartTime)
		if elapsed.Seconds() > 0 {
			method.RequestsPerSec = float64(method.TotalRequests) / elapsed.Seconds()
		}

		method.PercentComplete = (elapsed.Seconds() / float64(duration)) * 100
		if method.PercentComplete > 100 {
			method.PercentComplete = 100
		}
	}
}

// DisplayProgress displays all progress bars
func (pm *ProgressManager) DisplayProgress() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	// Only clear and redraw if we've displayed before
	if !pm.firstDisplay {
		// Move up 3 lines and clear them
		fmt.Print("\033[3A")
		fmt.Print("\033[K\033[K\033[K")
	} else {
		pm.firstDisplay = false
	}

	// Display each method's progress
	methodNames := []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"}

	for _, methodName := range methodNames {
		if method, exists := pm.methods[methodName]; exists {
			filledChar, emptyChar, icon := getProgressBarStyle(methodName)

			const barWidth = 20
			progress := int(method.PercentComplete * float64(barWidth) / 100)
			progressBar := strings.Repeat(filledChar, progress) + strings.Repeat(emptyChar, barWidth-progress)

			elapsed := int(time.Since(method.StartTime).Seconds())

			fmt.Printf("    %s [%s] %s: %.1f%% | %ds/%ds | Requests: %d | RPS: %.1f\n",
				icon, progressBar, methodName, method.PercentComplete, elapsed, duration, method.TotalRequests, method.RequestsPerSec)
		} else {
			// Method not started yet
			_, emptyChar, icon := getProgressBarStyle(methodName)
			progressBar := strings.Repeat(emptyChar, 20)
			fmt.Printf("    %s [%s] %s: 0.0%% | 0s/%ds | Requests: 0 | RPS: 0.0\n",
				icon, progressBar, methodName, duration)
		}
	}
}

// StartProgressDisplay starts the progress display loop
func (pm *ProgressManager) StartProgressDisplay() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.DisplayProgress()
		case <-pm.stopChan:
			return
		}
	}
}

// Stop stops the progress display
func (pm *ProgressManager) Stop() {
	close(pm.stopChan)
}

// runallCmd represents the runall command
var runallCmd = &cobra.Command{
	Use:   "runall",
	Short: "Run comprehensive test suite with all RPC methods",
	Long: `Execute a comprehensive test suite that includes:

1. Configuration Generation: Creates test configuration with your API key
2. Data Directory Setup: Creates ./data/ directory for storing test files  
3. Account Seeding: Seeds 100 accounts from specified program using remote RPC
4. Method Testing: Runs all available RPC methods concurrently against target RPC
5. Progress Tracking: Real-time progress bars with live statistics
6. Comprehensive Results: Detailed performance metrics for each method and overall summary

Features:
• Dual RPC Architecture: Uses remote RPC (from config) for seeding, target RPC (--url) for testing
• Real-time Progress: Visual progress bars with completion percentage and live RPS
• Dynamic Latency: Automatic unit formatting (μs, ms, s) based on performance
• Performance Insights: Method comparison with fastest/slowest analysis
• Account Management: Automatic account rotation and batching optimization

Required Flags:
  --api-key: API key for remote RPC endpoint (saved to config for future use)
  --url: Target RPC endpoint URL for testing and benchmarking

Examples:
  # Basic comprehensive test
  rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com
  
  # Advanced test with custom settings
  rpc_test runall --api-key YOUR_API_KEY --url https://your-target-rpc.com --concurrency 10 --duration 30 --limit 200
  
  # Test against Lantern (common use case)
  rpc_test runall --api-key YOUR_FLUX_API_KEY --url http://localhost:8080`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting comprehensive RPC test suite...")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Step 1: Generate and save test configuration
		var config TestConfig

		//check if config.json exists
		if _, err := os.Stat("./config.json"); err == nil {
			fmt.Println("📋 Step 1: Loading existing test configuration...")
			showProgress("Loading config", 100)
			config, err = loadTestConfig("./config.json")
			if err != nil {
				log.Fatalf("Failed to load test config: %v", err)
			}
			showProgressComplete("Config loaded")
			fmt.Printf("✅ Configuration loaded successfully\n")
		} else {
			fmt.Println("📋 Step 1: Generating test configuration...")
			showProgress("Generating config", 100)
			configFile := "./config.json"
			if err := generateTestConfig(configFile); err != nil {
				log.Fatalf("Failed to generate test config: %v", err)
			}
			showProgressComplete("Config generated")
			fmt.Printf("✅ Test configuration saved to: %s\n", configFile)
			config, err = loadTestConfig(configFile)
			if err != nil {
				log.Fatalf("Failed to load test config: %v", err)
			}
			showProgressComplete("Config loaded")
			fmt.Printf("✅ Configuration loaded successfully\n")
		}

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

// getProgressBarStyle returns different progress bar styles for different methods
func getProgressBarStyle(methodName string) (string, string, string) {
	switch methodName {
	case "getAccountInfo":
		return "█", "░", "🔍" // Solid blocks with magnifying glass
	case "getMultipleAccounts":
		return "▓", "░", "📊" // Dark blocks with chart
	case "getProgramAccounts":
		return "▒", "░", "⚙️" // Medium blocks with gear
	default:
		return "█", "░", "⚡" // Default solid blocks with lightning
	}
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
	for _, program := range config.Programs {
		programID = program
		break
	}

	if programID == "" {
		return fmt.Errorf("no program found in default configuration")
	}

	// Use the config RPC URL for seeding (remote RPC)
	seedRPCURL := config.RemoteRPCURL

	fmt.Printf("  🔍 Using remote RPC for seeding: %s\n", config.RemoteRPCURL)
	fmt.Printf("  🔍 Fetching accounts from program %s...\n", programID[:8]+"...")

	// Create RPC client for seeding (using config RPC URL)
	rpcTest := methods.NewRPCTest(seedRPCURL, config.RPCAPIKey)

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

	// Create progress manager
	progressManager := NewProgressManager()

	// Register all methods
	for _, methodName := range methods {
		progressManager.RegisterMethod(methodName, duration)
	}

	// Start progress display in background
	go progressManager.StartProgressDisplay()

	// Give a moment for initial display and to avoid interference with starting messages
	time.Sleep(500 * time.Millisecond)

	var results []TestResult
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Run each method concurrently
	for i, methodName := range methods {
		wg.Add(1)
		go func(method string, methodIndex int) {
			defer wg.Done()

			result := runSingleMethod(method, accounts, methodIndex+1, len(methods), progressManager)

			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(methodName, i)
	}

	wg.Wait()

	// Stop progress display
	progressManager.Stop()

	// Wait for the display goroutine to finish
	time.Sleep(500 * time.Millisecond)

	// Simple completion message without complex clearing
	fmt.Println()
	fmt.Println("    ✅ All methods completed successfully!")
	fmt.Println()

	return results, nil
}

// runSingleMethod runs a single method test and returns the result
func runSingleMethod(methodName string, accounts []string, methodIndex, totalMethods int, progressManager *ProgressManager) TestResult {
	fmt.Printf("  🔄 [%d/%d] Starting %s test...\n", methodIndex, totalMethods, methodName)

	// Create RPC client with target RPC URL (from --url flag)
	rpcTest := methods.NewRPCTest(rpcURL, apiKey)

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

	// Progress update ticker
	progressTicker := time.NewTicker(500 * time.Millisecond)
	defer progressTicker.Stop()

	// Progress update goroutine
	go func() {
		for {
			select {
			case <-progressTicker.C:
				if time.Now().After(endTime) {
					return
				}
				mutex.Lock()
				progressManager.UpdateProgress(methodName, successCount, failureCount)
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
					var err error

					if methodName == "getMultipleAccounts" {
						// For getMultipleAccounts, use multiple accounts
						// Take up to 5 accounts for each request
						numAccounts := rand.Intn(10) + 5
						if len(accounts) < numAccounts {
							numAccounts = len(accounts)
						}

						// Create a batch of accounts starting from workerID
						var batchAccounts []string
						for i := 0; i < numAccounts; i++ {
							accountIndex := (workerID + i) % len(accounts)
							batchAccounts = append(batchAccounts, accounts[accountIndex])
						}

						err = Method(methodName, rpcTest, batchAccounts...)
					} else {
						// For other methods, use single account
						err = Method(methodName, rpcTest, accounts[workerID%len(accounts)])
					}

					reqDuration := time.Since(startReq)

					mutex.Lock()
					if err != nil {
						fmt.Printf("  ❌ Error: %v\n", err)
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

	// Final progress update
	progressManager.UpdateProgress(methodName, successCount, failureCount)

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
