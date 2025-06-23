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

	_ "rpc_test/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	RemoteRPCURL string                  `json:"rpc_url" binding:"required"`
	RPCAPIKey    string                  `json:"rpc_apikey"`
	Programs     []string                `json:"programs"`
	TargetRPCURL string                  `json:"target_rpc_url" binding:"required"`
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

// APIResponse represents a generic API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

var (
	testManager *TestManager
	serverPort  string
	serverHost  string
)

// @host localhost:8081
// @BasePath /

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
- GET /swagger/*any - Swagger documentation

Example:
  rpc_test server --port 8081 --host localhost`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting RPC Test Server with Gin...")
		fmt.Printf("📍 Server will be available at: http://%s:%s\n", serverHost, serverPort)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Initialize test manager
		testManager = &TestManager{
			tests: make(map[string]*RunningTest),
		}

		// Set Gin mode
		gin.SetMode(gin.ReleaseMode)

		// Create Gin router
		r := gin.Default()

		// Add CORS middleware
		r.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		})

		// Setup routes
		setupRoutes(r)

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
		fmt.Println("   GET /swagger/*any - Swagger documentation")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

// setupRoutes configures all the API routes
func setupRoutes(r *gin.Engine) {
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/")
	{

		api.GET("/", handleRoot)
		api.POST("/test", handleTest)
		api.GET("/tests", handleTests)
		api.GET("/test/:id", handleTestByID)
		api.DELETE("/test/:id", handleDeleteTest)
	}
}

// @Summary Get server information
// @Description Get information about the RPC Test Server
// @Tags info
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Router / [get]
func handleRoot(c *gin.Context) {
	response := APIResponse{
		Success: true,
		Message: "RPC Test Server is running",
		Data: map[string]interface{}{
			"service": "RPC Test Server",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"POST /test":        "Start a new test",
				"GET /test/{id}":    "Get test results",
				"GET /tests":        "List all tests",
				"DELETE /test/{id}": "Delete a test",
				"GET /swagger/*any": "Swagger documentation",
			},
			"available_methods": []string{"getAccountInfo", "getMultipleAccounts", "getProgramAccounts"},
		},
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Start a new test
// @Description Start a new RPC test with the provided configuration
// @Tags tests
// @Accept json
// @Produce json
// @Param test body TestRequest true "Test configuration"
// @Success 200 {object} TestResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /test [post]
func handleTest(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success:   false,
			Message:   fmt.Sprintf("Invalid request: %v", err),
			Timestamp: time.Now(),
		})
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
	response := TestResponse{
		Success:   true,
		Message:   "Test started successfully",
		TestID:    testID,
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary List all tests
// @Description Get a list of all tests (running, completed, or failed)
// @Tags tests
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Router /tests [get]
func handleTests(c *gin.Context) {
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

	response := APIResponse{
		Success: true,
		Message: "Tests retrieved successfully",
		Data: map[string]interface{}{
			"tests": tests,
			"count": len(tests),
		},
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Get test results
// @Description Get results for a specific test by ID
// @Tags tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} TestResponse
// @Failure 404 {object} APIResponse
// @Router /test/{id} [get]
func handleTestByID(c *gin.Context) {
	testID := c.Param("id")

	testManager.mutex.RLock()
	test, exists := testManager.tests[testID]
	testManager.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Message:   "Test not found",
			Timestamp: time.Now(),
		})
		return
	}

	if test.Results != nil {
		c.JSON(http.StatusOK, test.Results)
	} else {
		// Test still running, return status
		response := map[string]interface{}{
			"id":         testID,
			"status":     test.Status,
			"start_time": test.StartTime,
			"config":     test.Config,
		}
		c.JSON(http.StatusOK, response)
	}
}

// @Summary Delete a test
// @Description Delete a specific test by ID
// @Tags tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /test/{id} [delete]
func handleDeleteTest(c *gin.Context) {
	testID := c.Param("id")

	testManager.mutex.Lock()
	_, exists := testManager.tests[testID]
	if exists {
		delete(testManager.tests, testID)
	}
	testManager.mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Message:   "Test not found",
			Timestamp: time.Now(),
		})
		return
	}

	response := APIResponse{
		Success:   true,
		Message:   "Test deleted successfully",
		Data:      map[string]string{"test_id": testID},
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
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
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "8081", "Server port")
	serverCmd.Flags().StringVarP(&serverHost, "host", "s", "localhost", "Server host")
}
