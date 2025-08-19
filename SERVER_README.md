# RPC Test Server

The RPC Test Server provides an HTTP API for running Solana RPC endpoint tests programmatically. This allows you to integrate RPC testing into your applications, CI/CD pipelines, or monitoring systems.

## Features

- **HTTP API**: RESTful endpoints for test management
- **FastHTTP Performance**: High-performance HTTP server using FastHTTP
- **CORS Support**: Cross-origin resource sharing enabled
- **Comprehensive Results**: Detailed performance metrics and statistics
- **Real-time Testing**: Synchronous test execution with immediate results
- **Method Testing**: Supports all three RPC methods (getAccountInfo, getMultipleAccounts, getProgramAccounts)

## Quick Start

### 1. Start the Server

```bash
# Start server on default port 8080
go run server.go

# The server will display startup information including available endpoints
```

### 2. Test the API

```bash
# Check server status
curl http://localhost:8080/

# Start a test with default configuration
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{}'

# Start a test with specific programs
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{
    "programs": ["2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"]
  }'
```

## API Endpoints

### GET /
Server information and available endpoints.

**Response:**
```json
{
  "success": true,
  "message": "RPC Test Server is running",
  "data": {
    "service": "RPC Test Server",
    "version": "1.0.0",
    "endpoints": {
      "GET /": "Server information",
      "POST /test": "Start a new test"
    },
    "available_methods": ["getAccountInfo", "getMultipleAccounts", "getProgramAccounts"]
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### POST /test
Start a new RPC test. The test runs synchronously and returns results immediately.

**Request Body:**
```json
{
  "programs": ["2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"]
}
```

**Note**: If no request body is provided or if the JSON parsing fails, the server will use default configuration with a default program.

**Default Configuration:**
- **Remote RPC URL**: Uses default RPC URL from server configuration
- **Target RPC URL**: Same as remote RPC URL
- **Programs**: `["2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"]` (default)
- **Concurrency**: 50 (per method)
- **Duration**: 15 seconds (per method)
- **Limit**: 50 accounts

**Response:**
```json
{
  "success": true,
  "message": "Test completed successfully",
  "results": [
    {
      "method_name": "getAccountInfo",
      "duration": 15.0,
      "total_requests": 750,
      "success_count": 745,
      "failure_count": 5,
      "success_rate": 99.33,
      "requests_per_sec": 49.67,
      "min_latency": 45.23,
      "max_latency": 125.67,
      "avg_latency": 78.45
    },
    {
      "method_name": "getMultipleAccounts",
      "duration": 15.0,
      "total_requests": 720,
      "success_count": 718,
      "failure_count": 2,
      "success_rate": 99.72,
      "requests_per_sec": 47.87,
      "min_latency": 52.11,
      "max_latency": 156.23,
      "avg_latency": 89.76
    },
    {
      "method_name": "getProgramAccounts", 
      "duration": 15.0,
      "total_requests": 680,
      "success_count": 675,
      "failure_count": 5,
      "success_rate": 99.26,
      "requests_per_sec": 45.33,
      "min_latency": 125.45,
      "max_latency": 456.78,
      "avg_latency": 234.56
    }
  ],
  "timestamp": "2024-01-01T12:00:00Z",
  "duration": 45000000000
}
```

## Server Configuration

The server uses the following default configuration:

```go
// Default server settings
serverHost = "localhost"
serverPort = "8080"
rpcURL = "https://api.mainnet-beta.solana.com"
concurrency = 50
duration = 15
limit = 50
```

## Architecture

### Core Components

1. **FastHTTP Server**: High-performance HTTP server
2. **CORS Middleware**: Enables cross-origin requests
3. **Test Manager**: Manages test execution
4. **Method Execution**: Runs all three RPC methods concurrently
5. **JSON Responses**: All responses use JSON format with Sonic for performance

### Test Execution Flow

1. **Request Processing**: Parse incoming JSON request or use defaults
2. **Configuration Setup**: Set up method configurations for all three RPC methods
3. **Account Seeding**: Fetch program accounts from the specified programs
4. **Concurrent Testing**: Run all three methods simultaneously
5. **Results Aggregation**: Collect and format test results
6. **Response**: Return comprehensive test results

### Methods Tested

The server automatically tests all three RPC methods:

1. **getAccountInfo**: Tests account information retrieval
2. **getMultipleAccounts**: Tests batch account retrieval (5-15 accounts per request)
3. **getProgramAccounts**: Tests program account enumeration

## Error Handling

### Common Error Responses

```json
{
  "success": false,
  "message": "Error description",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Troubleshooting

1. **Server won't start**: Check if port 8080 is available
2. **Connection refused**: Verify server is running with `go run server.go`
3. **JSON parsing errors**: Ensure request body is valid JSON
4. **Test failures**: Check RPC endpoint connectivity and program addresses

## Usage Examples

### Basic Health Check
```bash
curl http://localhost:8080/
```

### Run Default Test
```bash
curl -X POST http://localhost:8080/test
```

### Custom Program Test
```bash
curl -X POST http://localhost:8080/test \
  -H "Content-Type: application/json" \
  -d '{
    "programs": [
      "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
      "2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c"
    ]
  }'
```

### Integration with Lantern Configurator

This server is designed to work with the [Lantern configuration tool](https://configurator.fluxrpc.com/):

1. Start the server: `go run server.go`
2. Access the configurator at https://configurator.fluxrpc.com/
3. The configurator will automatically connect to your local server
4. Run tests directly from the web interface

## Performance Notes

- Uses FastHTTP for high-performance HTTP handling
- Supports concurrent testing of multiple RPC methods
- Optimized JSON marshaling with Sonic
- Real-time progress tracking during test execution
- Automatic cleanup of temporary test files 