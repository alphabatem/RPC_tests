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
	apiKey       string
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
	fmt.Printf("Number of accounts: %d\n", len(accounts))

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

	// Add progress reporting
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		fmt.Println("\nProgress:")
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
				const barWidth = 30
				progress := int(percentComplete * float64(barWidth) / 100)
				progressBar := strings.Repeat("█", progress) + strings.Repeat("░", barWidth-progress)

				fmt.Printf("\r[%s] %.1f%% | %ds/%ds | Requests: %d | RPS: %.1f",
					progressBar, percentComplete, int(elapsed.Seconds()), duration, currentTotal, currentRPS)
				mutex.Unlock()
			case <-stop:
				return
			}
		}
	}()

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

	// Improved results formatting with clearer visual separation
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 TEST RESULTS SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🕒 Duration:         %.2f seconds\n", totalDuration.Seconds())
	fmt.Printf("🔢 Total Requests:    %d\n", totalRequests)
	fmt.Printf("✅ Successful:        %d (%.2f%%)\n", successCount, successRate)
	fmt.Printf("❌ Failed:            %d (%.2f%%)\n", failureCount, 100-successRate)
	fmt.Printf("⚡ Requests/second:   %.2f\n", requestsPerSecond)

	// Add latency statistics
	if successCount > 0 {
		avgLatency := totalLatency / time.Duration(successCount)
		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("⏱️  LATENCY STATISTICS")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("Min: %.2f μs\n", float64(minLatency.Microseconds()))
		fmt.Printf("Max: %.2f μs\n", float64(maxLatency.Microseconds()))
		fmt.Printf("Avg: %.2f μs\n", float64(avgLatency.Microseconds()))
	}
}
