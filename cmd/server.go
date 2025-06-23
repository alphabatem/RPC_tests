package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"rpc_test/methods"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// MethodConfig represents configuration for a specific method
type MethodConfig struct {
	Concurrency int  `json:"concurrency"`
	Duration    int  `json:"duration"`
	Limit       int  `json:"limit"`
	Enabled     bool `json:"enabled"`
}

// TestRequest represents a test request from the API
type TestRequest struct {
	RemoteRPCURL string                  `json:"rpc_url"`
	RPCAPIKey    string                  `json:"rpc_apikey"`
	Programs     []string                `json:"programs"`
	TargetRPCURL string                  `json:"target_rpc_url"`
	Methods      map[string]MethodConfig `json:"methods"`
	GlobalConfig MethodConfig            `json:"global_config"`
}

// TestResponse represents the response from a test
type TestResponse struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	TestID    string        `json:"test_id,omitempty"`
	Results   []TestResult  `json:"results,omitempty"`
	Overall   OverallResult `json:"overall,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
}

// TestManager manages running tests
type TestManager struct {
	tests map[string]*RunningTest
	mutex sync.RWMutex
}

// RunningTest represents a test that's currently running
type RunningTest struct {
	ID        string
	Config    TestRequest
	Status    string // "running", "completed", "failed"
	Results   *TestResponse
	StartTime time.Time
	EndTime   time.Time
	Progress  chan TestProgress
}

// TestProgress represents progress updates during test execution
type TestProgress struct {
	MethodName      string  `json:"method_name"`
	PercentComplete float64 `json:"percent_complete"`
	Requests        int64   `json:"requests"`
	RPS             float64 `json:"rps"`
	SuccessRate     float64 `json:"success_rate"`
}

var (
	testManager *TestManager
	serverPort  string
	serverHost  string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server for RPC testing",
	Long: `Start an HTTP server that accepts test configurations and runs RPC tests via API endpoints.

The server provides the following endpoints:
- POST /test - Start a new test
- GET /test/{id} - Get test results
- GET /tests - List all tests
- DELETE /test/{id} - Delete a test

Example:
  rpc_test server --port 8080 --host localhost`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting RPC Test Server...")
		fmt.Printf("📍 Server will be available at: http://%s:%s\n", serverHost, serverPort)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Initialize test manager
		testManager = &TestManager{
			tests: make(map[string]*RunningTest),
		}

		// Setup routes
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/test", handleTest)
		http.HandleFunc("/tests", handleTests)
		http.HandleFunc("/test/", handleTestByID)

		// Start server
		addr := fmt.Sprintf("%s:%s", serverHost, serverPort)
		fmt.Printf("✅ Server started successfully!\n")
		fmt.Printf("📡 Listening on: %s\n", addr)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("📋 Available endpoints:")
		fmt.Println("   POST /test     - Start a new test")
		fmt.Println("   GET /test/{id} - Get test results")
		fmt.Println("   GET /tests     - List all tests")
		fmt.Println("   DELETE /test/{id} - Delete a test")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

// handleRoot handles the root endpoint
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"service": "RPC Test Server",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"POST /test":        "Start a new test",
			"GET /test/{id}":    "Get test results",
			"GET /tests":        "List all tests",
			"DELETE /test/{id}": "Delete a test",
		},
		"available_methods": []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"},
		"timestamp":         time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// handleTest handles POST /test for starting new tests
func handleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.TargetRPCURL == "" {
		http.Error(w, "target_rpc_url is required", http.StatusBadRequest)
		return
	}

	// Set defaults
	if len(req.Programs) == 0 {
		req.Programs = []string{"2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"}
	}

	// Initialize methods if not provided
	if req.Methods == nil {
		req.Methods = make(map[string]MethodConfig)
	}

	// Set default global config
	if req.GlobalConfig.Concurrency == 0 {
		req.GlobalConfig.Concurrency = 5
	}
	if req.GlobalConfig.Duration == 0 {
		req.GlobalConfig.Duration = 15
	}

	// Set defaults for each method if not specified
	availableMethods := []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"}
	for _, method := range availableMethods {
		if config, exists := req.Methods[method]; exists {
			// Use global defaults if method config is incomplete
			if config.Concurrency == 0 {
				config.Concurrency = req.GlobalConfig.Concurrency
			}
			if config.Duration == 0 {
				config.Duration = req.GlobalConfig.Duration
			}
			if config.Limit == 0 {
				config.Limit = req.GlobalConfig.Limit
			}
			if !config.Enabled {
				config.Enabled = true // Default to enabled
			}
			req.Methods[method] = config
		} else {
			// Create default config for method
			req.Methods[method] = MethodConfig{
				Concurrency: req.GlobalConfig.Concurrency,
				Duration:    req.GlobalConfig.Duration,
				Limit:       req.GlobalConfig.Limit,
				Enabled:     true,
			}
		}
	}

	// Generate test ID
	testID := generateTestID()

	// Create running test
	runningTest := &RunningTest{
		ID:        testID,
		Config:    req,
		Status:    "running",
		StartTime: time.Now(),
		Progress:  make(chan TestProgress, 100),
	}

	// Register test
	testManager.mutex.Lock()
	testManager.tests[testID] = runningTest
	testManager.mutex.Unlock()

	// Start test in background
	go runTestAsync(runningTest)

	// Return immediate response
	w.Header().Set("Content-Type", "application/json")
	response := TestResponse{
		Success:   true,
		Message:   "Test started successfully",
		TestID:    testID,
		Timestamp: time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// handleTests handles GET /tests for listing all tests
func handleTests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testManager.mutex.RLock()
	defer testManager.mutex.RUnlock()

	tests := make([]map[string]interface{}, 0)
	for id, test := range testManager.tests {
		testInfo := map[string]interface{}{
			"id":         id,
			"status":     test.Status,
			"start_time": test.StartTime,
			"end_time":   test.EndTime,
			"config":     test.Config,
		}
		if test.Results != nil {
			testInfo["duration"] = test.Results.Duration
		}
		tests = append(tests, testInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"tests":     tests,
		"count":     len(tests),
		"timestamp": time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// handleTestByID handles GET and DELETE /test/{id}
func handleTestByID(w http.ResponseWriter, r *http.Request) {
	// Extract test ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.NotFound(w, r)
		return
	}
	testID := pathParts[2]

	testManager.mutex.RLock()
	test, exists := testManager.tests[testID]
	testManager.mutex.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Return test results
		w.Header().Set("Content-Type", "application/json")
		if test.Results != nil {
			json.NewEncoder(w).Encode(test.Results)
		} else {
			// Test still running, return status
			response := map[string]interface{}{
				"id":         testID,
				"status":     test.Status,
				"start_time": test.StartTime,
				"config":     test.Config,
			}
			json.NewEncoder(w).Encode(response)
		}

	case http.MethodDelete:
		// Delete test
		testManager.mutex.Lock()
		delete(testManager.tests, testID)
		testManager.mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"success":   true,
			"message":   "Test deleted successfully",
			"test_id":   testID,
			"timestamp": time.Now(),
		}
		json.NewEncoder(w).Encode(response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// runTestAsync runs a test asynchronously
func runTestAsync(test *RunningTest) {
	defer func() {
		test.EndTime = time.Now()
		if test.Results != nil {
			test.Results.Duration = test.EndTime.Sub(test.StartTime)
		}
	}()

	// Create temporary config file
	configFile := fmt.Sprintf("./data/server_config_%s.json", test.ID)
	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		test.Status = "failed"
		test.Results = &TestResponse{
			Success:   false,
			Message:   fmt.Sprintf("Failed to create config directory: %v", err),
			TestID:    test.ID,
			Timestamp: time.Now(),
		}
		return
	}

	// Generate config
	config := TestConfig{
		RemoteRPCURL: test.Config.RemoteRPCURL,
		RPCAPIKey:    test.Config.RPCAPIKey,
		Programs:     test.Config.Programs,
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		test.Status = "failed"
		test.Results = &TestResponse{
			Success:   false,
			Message:   fmt.Sprintf("Failed to marshal config: %v", err),
			TestID:    test.ID,
			Timestamp: time.Now(),
		}
		return
	}

	if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
		test.Status = "failed"
		test.Results = &TestResponse{
			Success:   false,
			Message:   fmt.Sprintf("Failed to write config file: %v", err),
			TestID:    test.ID,
			Timestamp: time.Now(),
		}
		return
	}

	// Seed accounts
	accountsFile := fmt.Sprintf("./data/server_accounts_%s.txt", test.ID)
	if err := seedAccountsFromProgram(accountsFile, config); err != nil {
		test.Status = "failed"
		test.Results = &TestResponse{
			Success:   false,
			Message:   fmt.Sprintf("Failed to seed accounts: %v", err),
			TestID:    test.ID,
			Timestamp: time.Now(),
		}
		return
	}

	// Run tests for each enabled method
	var allResults []TestResult
	var wg sync.WaitGroup
	var resultsMutex sync.Mutex

	// Store original global values
	originalRPCURL := rpcURL
	originalConcurrency := concurrency
	originalDuration := duration
	originalLimit := limit

	// Set target RPC URL
	rpcURL = test.Config.TargetRPCURL

	// Run each enabled method
	for methodName, methodConfig := range test.Config.Methods {
		if !methodConfig.Enabled {
			continue
		}

		wg.Add(1)
		go func(method string, config MethodConfig) {
			defer wg.Done()

			// Set method-specific configuration
			concurrency = config.Concurrency
			duration = config.Duration
			limit = config.Limit

			// Run the method test
			result := runServerMethod(method, accountsFile, &test.Config)

			// Store result
			resultsMutex.Lock()
			allResults = append(allResults, result)
			resultsMutex.Unlock()
		}(methodName, methodConfig)
	}

	// Wait for all methods to complete
	wg.Wait()

	// Restore original values
	rpcURL = originalRPCURL
	concurrency = originalConcurrency
	duration = originalDuration
	limit = originalLimit

	// Check if we have any results
	if len(allResults) == 0 {
		test.Status = "failed"
		test.Results = &TestResponse{
			Success:   false,
			Message:   "No methods were enabled or all methods failed",
			TestID:    test.ID,
			Timestamp: time.Now(),
		}
		return
	}

	// Calculate overall results
	overall := calculateOverallResults(allResults)

	test.Status = "completed"
	test.Results = &TestResponse{
		Success:   true,
		Message:   "Test completed successfully",
		TestID:    test.ID,
		Results:   allResults,
		Overall:   overall,
		Timestamp: time.Now(),
	}

	// Cleanup temporary files
	os.Remove(configFile)
	os.Remove(accountsFile)
}

// runServerMethod runs a single method test with the given configuration
func runServerMethod(methodName string, accountsFile string, testConfig *TestRequest) TestResult {
	// Load accounts from file
	data, err := os.ReadFile(accountsFile)
	if err != nil {
		return TestResult{
			MethodName:     methodName,
			Duration:       0,
			TotalRequests:  0,
			SuccessCount:   0,
			FailureCount:   1,
			RequestsPerSec: 0,
			SuccessRate:    0,
		}
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
		return TestResult{
			MethodName:     methodName,
			Duration:       0,
			TotalRequests:  0,
			SuccessCount:   0,
			FailureCount:   1,
			RequestsPerSec: 0,
			SuccessRate:    0,
		}
	}

	// Get method configuration
	methodConfig := testConfig.Methods[methodName]

	// Apply limit if specified
	if methodConfig.Limit > 0 && methodConfig.Limit < len(accounts) {
		accounts = accounts[:methodConfig.Limit]
	}

	// Create RPC client
	rpcTest := methods.NewRPCTest(rpcURL)

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(methodConfig.Duration) * time.Second)

	var wg sync.WaitGroup
	var successCount, failureCount int64
	var mutex sync.Mutex

	// Create channels for workers
	stop := make(chan struct{})

	// Collect statistics
	var totalLatency time.Duration
	var minLatency time.Duration = time.Hour
	var maxLatency time.Duration

	// Start workers
	for i := 0; i < methodConfig.Concurrency; i++ {
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
						numAccounts := rand.Intn(10) + 5
						if len(accounts) < numAccounts {
							numAccounts = len(accounts)
						}
						var batchAccounts []string
						for i := 0; i < numAccounts; i++ {
							accountIndex := (workerID + i) % len(accounts)
							batchAccounts = append(batchAccounts, accounts[accountIndex])
						}
						err = Method(methodName, rpcTest, batchAccounts...)
					} else {
						err = Method(methodName, rpcTest, accounts[workerID%len(accounts)])
					}
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
	time.Sleep(time.Duration(methodConfig.Duration) * time.Second)
	close(stop)

	// Wait for all workers to finish
	wg.Wait()

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

// generateTestID generates a unique test ID
func generateTestID() string {
	return fmt.Sprintf("test_%d", time.Now().UnixNano())
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Add server-specific flags
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "8080", "Server port")
	serverCmd.Flags().StringVarP(&serverHost, "host", "h", "localhost", "Server host")
}
