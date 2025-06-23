# RPC Test Server

The RPC Test Server provides an HTTP API for running Solana RPC endpoint tests programmatically. This allows you to integrate RPC testing into your applications, CI/CD pipelines, or monitoring systems.

## Features

- **HTTP API**: RESTful endpoints for test management
- **Asynchronous Testing**: Tests run in the background, allowing non-blocking API calls
- **Method-Specific Configuration**: Configure different settings for each RPC method
- **Test Management**: Start, monitor, and retrieve test results
- **Multiple Test Support**: Run and manage multiple tests simultaneously
- **Comprehensive Results**: Detailed performance metrics and statistics
- **Cleanup**: Automatic cleanup of temporary files

## Quick Start

### 1. Start the Server

```bash
# Start server on default port 8080
./rpc_test server

# Start server on custom port and host
./rpc_test server --port 9000 --host 0.0.0.0
```

### 2. Test the API

```bash
# Check server status
curl http://localhost:8080/

# Start a test
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{
    "target_rpc_url": "https://api.mainnet-beta.solana.com",
    "concurrency": 5,
    "duration": 15
  }'
```

## API Endpoints

### GET /
Server information and available endpoints.

**Response:**
```json
{
  "service": "RPC Test Server",
  "version": "1.0.0",
  "endpoints": {
    "POST /test": "Start a new test",
    "GET /test/{id}": "Get test results",
    "GET /tests": "List all tests",
    "DELETE /test/{id}": "Delete a test"
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### POST /test
Start a new RPC test.

**Request Body:**
```json
{
  "remote_rpc_url": "https://us.rpc.fluxbeam.xyz",
  "rpc_apikey": "YOUR_API_KEY_HERE",
  "programs": ["2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"],
  "target_rpc_url": "https://api.mainnet-beta.solana.com",
  "global_config": {
    "concurrency": 5,
    "duration": 15,
    "limit": 0
  },
  "methods": {
    "getAccountInfo": {
      "concurrency": 10,
      "duration": 20,
      "limit": 50,
      "enabled": true
    },
    "getMultipleAccounts": {
      "concurrency": 5,
      "duration": 15,
      "limit": 100,
      "enabled": true
    },
    "getProgramAccounts": {
      "concurrency": 3,
      "duration": 10,
      "limit": 25,
      "enabled": false
    }
  }
}
```

**Required Fields:**
- `target_rpc_url`: The RPC endpoint to test

**Optional Fields:**
- `remote_rpc_url`: RPC endpoint for seeding accounts (default: Fluxbeam)
- `rpc_apikey`: API key for remote RPC (default: none)
- `programs`: Array of program IDs to seed accounts from (default: Fluxbeam program)
- `global_config`: Global configuration defaults for methods
- `methods`: Method-specific configurations

**Global Config Options:**
- `concurrency`: Default number of concurrent requests (default: 5)
- `duration`: Default test duration in seconds (default: 15)
- `limit`: Default limit for number of accounts (default: 0 = no limit)

**Method Config Options:**
- `concurrency`: Number of concurrent requests for this method
- `duration`: Test duration in seconds for this method
- `limit`: Limit number of accounts for this method
- `enabled`: Whether to run this method (default: true)

**Response:**
```json
{
  "success": true,
  "message": "Test started successfully",
  "test_id": "test_1704067200000000000",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### GET /test/{id}
Get test results or status.

**Response (Running Test):**
```json
{
  "id": "test_1704067200000000000",
  "status": "running",
  "start_time": "2024-01-01T12:00:00Z",
  "config": {
    "target_rpc_url": "https://api.mainnet-beta.solana.com",
    "concurrency": 5,
    "duration": 15
  }
}
```

**Response (Completed Test):**
```json
{
  "success": true,
  "message": "Test completed successfully",
  "test_id": "test_1704067200000000000",
  "results": [
    {
      "method_name": "getAccountInfo",
      "duration": "15.2s",
      "total_requests": 150,
      "success_count": 148,
      "failure_count": 2,
      "requests_per_sec": 9.87,
      "success_rate": 98.67,
      "min_latency": "45ms",
      "max_latency": "1.2s",
      "avg_latency": "120ms"
    }
  ],
  "overall": {
    "total_duration": "15.2s",
    "total_requests": 450,
    "total_success": 445,
    "total_failure": 5,
    "overall_rps": 29.61,
    "overall_success_rate": 98.89
  },
  "timestamp": "2024-01-01T12:00:15Z",
  "duration": "15.2s"
}
```

### GET /tests
List all tests.

**Response:**
```json
{
  "tests": [
    {
      "id": "test_1704067200000000000",
      "status": "completed",
      "start_time": "2024-01-01T12:00:00Z",
      "end_time": "2024-01-01T12:00:15Z",
      "config": {
        "target_rpc_url": "https://api.mainnet-beta.solana.com",
        "concurrency": 5,
        "duration": 15
      },
      "duration": "15.2s"
    }
  ],
  "count": 1,
  "timestamp": "2024-01-01T12:00:20Z"
}
```

### DELETE /test/{id}
Delete a test.

**Response:**
```json
{
  "success": true,
  "message": "Test deleted successfully",
  "test_id": "test_1704067200000000000",
  "timestamp": "2024-01-01T12:00:25Z"
}
```

## Configuration Examples

### Simple Configuration (Global Defaults)
```json
{
  "target_rpc_url": "https://api.mainnet-beta.solana.com",
  "global_config": {
    "concurrency": 5,
    "duration": 15
  }
}
```

### Method-Specific Configuration
```json
{
  "target_rpc_url": "https://api.mainnet-beta.solana.com",
  "global_config": {
    "concurrency": 3,
    "duration": 10
  },
  "methods": {
    "getAccountInfo": {
      "concurrency": 10,
      "duration": 20,
      "enabled": true
    },
    "getMultipleAccounts": {
      "concurrency": 5,
      "duration": 15,
      "enabled": true
    },
    "getProgramAccounts": {
      "enabled": false
    }
  }
}
```

### Advanced Configuration with Limits
```json
{
  "remote_rpc_url": "https://us.rpc.fluxbeam.xyz",
  "rpc_apikey": "YOUR_API_KEY_HERE",
  "programs": ["2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"],
  "target_rpc_url": "https://api.mainnet-beta.solana.com",
  "global_config": {
    "concurrency": 5,
    "duration": 15,
    "limit": 100
  },
  "methods": {
    "getAccountInfo": {
      "concurrency": 15,
      "duration": 30,
      "limit": 200,
      "enabled": true
    },
    "getMultipleAccounts": {
      "concurrency": 8,
      "duration": 20,
      "limit": 150,
      "enabled": true
    },
    "getProgramAccounts": {
      "concurrency": 3,
      "duration": 10,
      "limit": 50,
      "enabled": true
    }
  }
}
```

## Client Examples

### Python Client

See `examples/client_example.py` for a complete Python client implementation.

```python
import requests

# Start a test
response = requests.post('http://localhost:8080/test', json={
    'target_rpc_url': 'https://api.mainnet-beta.solana.com',
    'concurrency': 5,
    'duration': 15
})

test_id = response.json()['test_id']

# Get results
results = requests.get(f'http://localhost:8080/test/{test_id}').json()
```

### cURL Examples

See `examples/curl_examples.sh` for complete cURL examples.

```bash
# Start test
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{"target_rpc_url": "https://api.mainnet-beta.solana.com"}'

# Get results
curl http://localhost:8080/test/test_1704067200000000000
```

## Integration Examples

### CI/CD Pipeline

```yaml
# GitHub Actions example
- name: Test RPC Endpoint
  run: |
    # Start test
    response=$(curl -s -X POST http://localhost:8080/test \
      -H "Content-Type: application/json" \
      -d '{"target_rpc_url": "${{ secrets.RPC_URL }}"}')
    
    test_id=$(echo $response | jq -r '.test_id')
    
    # Wait for completion
    while true; do
      status=$(curl -s http://localhost:8080/test/$test_id | jq -r '.status')
      if [ "$status" = "completed" ]; then
        break
      elif [ "$status" = "failed" ]; then
        exit 1
      fi
      sleep 5
    done
    
    # Check results
    results=$(curl -s http://localhost:8080/test/$test_id)
    success_rate=$(echo $results | jq -r '.overall.overall_success_rate')
    
    if (( $(echo "$success_rate < 95" | bc -l) )); then
      echo "Success rate too low: $success_rate%"
      exit 1
    fi
```

### Monitoring Dashboard

```javascript
// JavaScript example for monitoring dashboard
async function runRPCTest() {
    // Start test
    const startResponse = await fetch('/test', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            target_rpc_url: 'https://api.mainnet-beta.solana.com',
            concurrency: 10,
            duration: 30
        })
    });
    
    const { test_id } = await startResponse.json();
    
    // Poll for results
    const pollInterval = setInterval(async () => {
        const response = await fetch(`/test/${test_id}`);
        const result = await response.json();
        
        if (result.status === 'completed') {
            clearInterval(pollInterval);
            updateDashboard(result);
        }
    }, 5000);
}
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200`: Success
- `400`: Bad Request (invalid JSON or missing required fields)
- `404`: Test not found
- `405`: Method not allowed
- `500`: Internal server error

Error responses include a message explaining the issue:

```json
{
  "error": "target_rpc_url is required"
}
```

## Configuration

### Server Configuration

```bash
# Default configuration
./rpc_test server

# Custom configuration
./rpc_test server --port 9000 --host 0.0.0.0
```

### Environment Variables

You can also set environment variables:

```bash
export RPC_TEST_SERVER_PORT=9000
export RPC_TEST_SERVER_HOST=0.0.0.0
./rpc_test server
```

## Security Considerations

1. **API Keys**: Store API keys securely and don't expose them in client-side code
2. **Network Access**: Consider firewall rules to restrict server access
3. **Rate Limiting**: Implement rate limiting for production use
4. **Authentication**: Add authentication for production deployments

## Troubleshooting

### Common Issues

1. **Server won't start**: Check if port is already in use
2. **Tests fail**: Verify RPC endpoints are accessible
3. **Memory issues**: Monitor server memory usage with many concurrent tests
4. **Network timeouts**: Adjust timeout settings for slow RPC endpoints

### Debug Mode

Enable debug logging by setting the log level:

```bash
export RPC_TEST_LOG_LEVEL=debug
./rpc_test server
```

## Performance

- **Concurrent Tests**: Server can handle multiple tests simultaneously
- **Memory Usage**: Each test uses temporary files that are cleaned up automatically
- **Network**: Tests use the same network configuration as CLI commands
- **Scaling**: For high load, consider running multiple server instances behind a load balancer

## Contributing

To extend the server functionality:

1. Add new endpoints in `cmd/server.go`
2. Update the test execution logic in `runTestAsync()`
3. Add new configuration options as needed
4. Update documentation and examples 