package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"rpc_test/methods"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
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

type TestRequestSimple struct {
	Programs []string `json:"programs,omitempty"`
}

// TestRequest represents a test request from the API
type TestRequest struct {
	RemoteRPCURL string                  `json:"rpc_url,omitempty"`
	Programs     []string                `json:"programs,omitempty"`
	TargetRPCURL string                  `json:"target_rpc_url,omitempty"`
	Methods      map[string]MethodConfig `json:"methods,omitempty"`
	GlobalConfig MethodConfig            `json:"global_config,omitempty"`
}

// TestResponse represents the response from a test
type TestResponse struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	TestID    string        `json:"test_id,omitempty"`
	Results   []TestResult  `json:"results,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
}

// TestResult represents the result of a single method test
type TestResult struct {
	MethodName       string  `json:"method_name"`
	Duration         int64   `json:"duration_micros"`
	TotalRequests    int64   `json:"total_requests"`
	SuccessCount     int64   `json:"success_count"`
	FailureCount     int64   `json:"failure_count"`
	RequestsPerSec   float64 `json:"requests_per_sec"`
	SuccessRate      float64 `json:"success_rate"`
	MinLatencyMicros int64   `json:"min_latency_micros"`
	MaxLatencyMicros int64   `json:"max_latency_micros"`
	AvgLatencyMicros int64   `json:"avg_latency_micros"`
}

// TestConfig represents the configuration for seeding
type TestConfig struct {
	RemoteRPCURL string   `json:"rpc_url"`
	RPCAPIKey    string   `json:"rpc_apikey"`
	Programs     []string `json:"programs"`
}

// TestManager manages running tests
type TestManager struct {
	tests map[string]*RunningTest
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

// APIResponse represents a generic API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

var (
	testManager *TestManager
	serverPort  = "8888"
	serverHost  = "localhost"

	// Global variables for RPC testing
	rpcURL      = "http://localhost:8080"
	concurrency = 1
	duration    = 5
	limit       = 50
)

// JSON response helper
func writeJSONResponse(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)

	jsonData, err := sonic.Marshal(data)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(`{"success":false,"message":"JSON marshal error","timestamp":"` + strconv.FormatInt(time.Now().Unix(), 10) + `"}`)
		return
	}

	ctx.Write(jsonData)
}

// CORS middleware
func corsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		next(ctx)
	}
}

func main() {
	fmt.Println("🚀 Starting RPC Test Server with FastHTTP...")
	fmt.Printf("📍 Local access: http://localhost:%s\n", serverPort)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Initialize test manager
	testManager = &TestManager{
		tests: make(map[string]*RunningTest),
	}

	// Create router
	r := router.New()

	// Setup routes
	setupRoutes(r)

	// Start server
	addr := fmt.Sprintf("%s:%s", serverHost, serverPort)
	fmt.Printf("✅ Server started successfully!\n")
	fmt.Printf("📡 Listening on: %s\n", addr)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 Available endpoints:")
	fmt.Println("   GET /          - Server information")
	fmt.Println("   POST /test     - Start a new test")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := fasthttp.ListenAndServe(addr, corsMiddleware(r.Handler)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes configures all the API routes
func setupRoutes(r *router.Router) {
	// API routes
	r.GET("/", handleRoot)
	r.POST("/test", handleTest)
}

func handleRoot(ctx *fasthttp.RequestCtx) {
	response := APIResponse{
		Success: true,
		Message: "RPC Test Server is running",
		Data: map[string]interface{}{
			"service": "RPC Test Server",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"GET /":      "Server information",
				"POST /test": "Start a new test",
			},
			"available_methods": []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"},
		},
		Timestamp: time.Now(),
	}

	writeJSONResponse(ctx, fasthttp.StatusOK, response)
}

func handleTest(ctx *fasthttp.RequestCtx) {
	var reqBody TestRequestSimple
	var req TestRequest
	if err := json.Unmarshal(ctx.PostBody(), &reqBody); err != nil {
		req = TestRequest{
			RemoteRPCURL: rpcURL,
			TargetRPCURL: rpcURL,
			Programs:     []string{"2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"},
			GlobalConfig: MethodConfig{
				Concurrency: concurrency,
				Duration:    duration,
				Limit:       limit,
				Enabled:     true,
			},
		}
	} else {
		req = TestRequest{
			RemoteRPCURL: rpcURL,
			TargetRPCURL: rpcURL,
			Programs:     reqBody.Programs,
			GlobalConfig: MethodConfig{
				Concurrency: concurrency,
				Duration:    duration,
				Limit:       limit,
				Enabled:     true,
			},
		}
	}

	fmt.Println(req)
	if req.Methods == nil {
		req.Methods = make(map[string]MethodConfig)
	}

	// Set defaults for each method if not specified
	for _, method := range []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"} {
		req.Methods[method] = MethodConfig{
			Concurrency: req.GlobalConfig.Concurrency,
			Duration:    req.GlobalConfig.Duration,
			Limit:       req.GlobalConfig.Limit,
			Enabled:     true,
		}
	}

	// Create running test
	runningTest := &RunningTest{
		Config:    req,
		Status:    "running",
		StartTime: time.Now(),
		Progress:  make(chan TestProgress, 100),
	}

	// change to running test and get data
	response := runTestAsync(runningTest)

	writeJSONResponse(ctx, fasthttp.StatusOK, response)
}

// Method executes a specific RPC method
func Method(name string, rpcTest *methods.RPCTest, account ...string) error {
	switch name {
	case "getAccountInfo":
		return rpcTest.GetAccountInfo(account[0])
	case "getMultipleAccounts":
		return rpcTest.GetMultipleAccounts(account...)
	case "getProgramAccounts":
		return rpcTest.GetProgramAccounts(account[0])
	default:
		return fmt.Errorf("invalid method: %s", name)
	}
}

// runTestAsync runs a test synchronously
func runTestAsync(test *RunningTest) *TestResponse {
	defer func() {
		test.EndTime = time.Now()
		if test.Results != nil {
			test.Results.Duration = test.EndTime.Sub(test.StartTime)
		}
	}()

	// Run tests for each enabled method
	var allResults []TestResult

	// Store original global values
	originalRPCURL := rpcURL
	originalConcurrency := concurrency
	originalDuration := duration
	originalLimit := limit

	// Set target RPC URL
	rpcURL = test.Config.TargetRPCURL
	var accounts []string
	var err error

	accounts, err = loadAccountsFromFile("./data/test_accounts.txt", test.Config)
	if err != nil {
		fmt.Println("Error loading accounts:", err)
		test.Status = "failed"
		test.Results = &TestResponse{
			Success:   false,
			Message:   fmt.Sprintf("Failed to load accounts: %v", err),
			TestID:    test.ID,
			Timestamp: time.Now(),
		}
		return test.Results
	}

	// Run methods in a specific order
	methodOrder := []string{"getProgramAccounts", "getAccountInfo", "getMultipleAccounts"}

	for _, methodName := range methodOrder {
		methodConfig, exists := test.Config.Methods[methodName]
		if !exists || !methodConfig.Enabled {
			continue
		}

		// Set method-specific configuration
		concurrency = methodConfig.Concurrency
		duration = methodConfig.Duration
		limit = methodConfig.Limit

		// Run the method test
		result := runServerMethod(methodName, &test.Config, accounts)
		allResults = append(allResults, result)

		fmt.Printf("Completed %s: %d requests in %v\n",
			methodName, result.TotalRequests, time.Duration(result.Duration)*time.Microsecond)
	}

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
		return test.Results
	}

	test.Status = "completed"
	test.Results = &TestResponse{
		Success:   true,
		Message:   "Test completed successfully",
		TestID:    test.ID,
		Results:   allResults,
		Timestamp: time.Now(),
	}
	return test.Results
}

// runServerMethod runs a single method test with the given configuration
func runServerMethod(methodName string, testConfig *TestRequest, accounts []string) TestResult {
	if len(accounts) == 0 {
		return TestResult{
			MethodName:       methodName,
			Duration:         0,
			TotalRequests:    0,
			SuccessCount:     0,
			FailureCount:     1,
			RequestsPerSec:   0,
			SuccessRate:      0,
			MinLatencyMicros: 0,
			MaxLatencyMicros: 0,
			AvgLatencyMicros: 0,
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

	var successCount, failureCount int64
	var totalLatency time.Duration
	var minLatency time.Duration = time.Hour
	var maxLatency time.Duration

	// Run test synchronously for the duration
	accountIndex := 0
	for time.Now().Before(endTime) {
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
				idx := (accountIndex + i) % len(accounts)
				batchAccounts = append(batchAccounts, accounts[idx])
			}
			err = Method(methodName, rpcTest, batchAccounts...)
		} else if methodName == "getProgramAccounts" {
			err = Method(methodName, rpcTest, testConfig.Programs...)
		} else {
			err = Method(methodName, rpcTest, accounts[accountIndex%len(accounts)])
		}

		reqDuration := time.Since(startReq)
		accountIndex++

		if err != nil {
			failureCount++
			fmt.Printf("Error in %s: %v\n", methodName, err)
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
	}

	// Calculate results
	totalDuration := time.Since(startTime)
	totalRequests := successCount + failureCount
	requestsPerSecond := float64(totalRequests) / totalDuration.Seconds()
	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(successCount) / float64(totalRequests) * 100
	}

	var avgLatency time.Duration
	if successCount > 0 {
		avgLatency = totalLatency / time.Duration(successCount)
	}

	return TestResult{
		MethodName:       methodName,
		Duration:         totalDuration.Microseconds(),
		TotalRequests:    totalRequests,
		SuccessCount:     successCount,
		FailureCount:     failureCount,
		RequestsPerSec:   requestsPerSecond,
		SuccessRate:      successRate,
		MinLatencyMicros: minLatency.Microseconds(),
		MaxLatencyMicros: maxLatency.Microseconds(),
		AvgLatencyMicros: avgLatency.Microseconds(),
	}
}

// Load accounts from file
func loadAccountsFromFile(accountsFile string, testConfig TestRequest) ([]string, error) {
	if testConfig.Programs[0] != "2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c" {
		newFile := fmt.Sprintf("./data/test_accounts_%s.txt", testConfig.Programs[0])
		defer os.Remove(newFile)
		err := seedAccountsFromProgram(newFile, TestConfig{
			RemoteRPCURL: rpcURL,
			Programs:     testConfig.Programs,
		})
		if err != nil {
			return nil, err
		}
		accountsFile = newFile
	}

	data, err := os.ReadFile(accountsFile)
	if err != nil || len(data) == 0 {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var accounts []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			accounts = append(accounts, line)
		}
	}
	return accounts, nil
}

// seedAccountsFromProgram seeds accounts from a program
func seedAccountsFromProgram(accountsFile string, config TestConfig) error {
	// Create RPC client for seeding
	rpcTest := methods.NewRPCTest(config.RemoteRPCURL)

	// Seed from the first program (or use default)
	programAddress := "2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"
	if len(config.Programs) > 0 {
		programAddress = config.Programs[0]
	}

	// Use a reasonable limit for seeding
	seedLimit := 100
	if limit > 0 {
		seedLimit = limit
	}

	return rpcTest.SeedProgramAccounts(programAddress, accountsFile, seedLimit)
}

// generateTestID generates a unique test ID
func generateTestID() string {
	return fmt.Sprintf("test_%d", time.Now().UnixNano())
}
