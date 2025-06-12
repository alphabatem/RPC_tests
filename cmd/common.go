package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"rpc_test/methods"
)

// Common variables for all commands
var (
	rpcURL       string
	concurrency  int
	duration     int
	accounts     []string
	accountsFile string
	limit        int
)

func Method(name string, rpcTest *methods.RPCTest, account string) error {
	switch name {
	case "getAccountInfo":
		return rpcTest.GetAccountInfo(account)
	case "getMultipleAccounts":
		return rpcTest.GetMultipleAccounts(account)
	case "getProgramAccounts":
		return rpcTest.GetProgramAccounts(account)
	default:
		return fmt.Errorf("invalid method: %s", name)
	}
}

// RunMethodTest runs a performance test for a specific RPC method
func RunMethodTest(methodName string) {
	// Load accounts from file if provided
	if accountsFile != "" {
		data, err := os.ReadFile(accountsFile)
		if err != nil {
			log.Fatalf("Failed to read accounts file: %v", err)
		}
		// Parse accounts from file
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				accounts = append(accounts, line)
			}
		}
	}

	if len(accounts) == 0 {
		log.Fatalf("No accounts provided. Use --account or --account-file to specify accounts")
	}

	// Apply limit if specified
	totalAccounts := len(accounts)
	if limit > 0 && limit < totalAccounts {
		accounts = accounts[:limit]
		fmt.Printf("Limiting to %d accounts out of %d available\n", limit, totalAccounts)
	}

	// Create RPC client
	rpcTest := methods.NewRPCTest(rpcURL)

	// Run the stress test
	fmt.Printf("Starting %s test with %d concurrent requests for %d seconds\n",
		methodName, concurrency, duration)
	fmt.Printf("RPC URL: %s\n", rpcURL)
	fmt.Printf("Accounts: %v\n", accounts)

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

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

	// Calculate and display results
	totalDuration := time.Since(startTime)
	totalRequests := successCount + failureCount
	requestsPerSecond := float64(totalRequests) / totalDuration.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100

	fmt.Println("\nTest Results:")
	fmt.Printf("Total Duration: %.2f seconds\n", totalDuration.Seconds())
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful Requests: %d (%.2f%%)\n", successCount, successRate)
	fmt.Printf("Failed Requests: %d (%.2f%%)\n", failureCount, 100-successRate)
	fmt.Printf("Requests per second: %.2f\n", requestsPerSecond)

	// Add latency statistics
	if successCount > 0 {
		avgLatency := totalLatency / time.Duration(successCount)
		fmt.Printf("\nLatency Statistics:\n")
		fmt.Printf("Min: %.2f ms\n", float64(minLatency.Microseconds())/1000)
		fmt.Printf("Max: %.2f ms\n", float64(maxLatency.Microseconds())/1000)
		fmt.Printf("Avg: %.2f ms\n", float64(avgLatency.Microseconds())/1000)
	}
}
