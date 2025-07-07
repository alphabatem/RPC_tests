#!/bin/bash

# RPC Test Server API Examples using curl
# Make sure the server is running: ./rpc_test server

SERVER_URL="http://localhost:8080"

echo "🔧 RPC Test Server API Examples"
echo "=================================="

# Function to print response
print_response() {
    echo "Status: $1"
    echo "Response:"
    echo "$2" | jq '.' 2>/dev/null || echo "$2"
    echo "----------------------------------"
}

# 1. Check server status
echo "1. Checking server status..."
response=$(curl -s -w "%{http_code}" "$SERVER_URL/")
status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

# 2. Start a test with method-specific configurations
echo "2. Starting a test with method-specific configurations..."
test_config='{
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
}'

response=$(curl -s -w "%{http_code}" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$test_config" \
  "$SERVER_URL/test")

status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

# Extract test ID from response
test_id=$(echo "$body" | jq -r '.test_id' 2>/dev/null)
if [ "$test_id" = "null" ] || [ -z "$test_id" ]; then
    echo "❌ Failed to get test ID"
    exit 1
fi

echo "📋 Test ID: $test_id"

# 3. Start a simple test using global defaults
echo "3. Starting a simple test using global defaults..."
simple_test_config='{
  "target_rpc_url": "https://api.mainnet-beta.solana.com",
  "global_config": {
    "concurrency": 3,
    "duration": 10
  }
}'

response=$(curl -s -w "%{http_code}" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$simple_test_config" \
  "$SERVER_URL/test")

status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

# Extract second test ID
test_id2=$(echo "$body" | jq -r '.test_id' 2>/dev/null)
if [ "$test_id2" != "null" ] && [ -n "$test_id2" ]; then
    echo "📋 Simple Test ID: $test_id2"
fi

# 4. Check test status
echo "4. Checking test status..."
response=$(curl -s -w "%{http_code}" "$SERVER_URL/test/$test_id")
status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

# 5. List all tests
echo "5. Listing all tests..."
response=$(curl -s -w "%{http_code}" "$SERVER_URL/tests")
status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

# 6. Wait for test completion (polling)
echo "6. Waiting for test completion..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    response=$(curl -s "$SERVER_URL/test/$test_id")
    status=$(echo "$response" | jq -r '.status' 2>/dev/null)
    
    echo "Attempt $((attempt + 1))/$max_attempts - Status: $status"
    
    if [ "$status" = "completed" ]; then
        echo "✅ Test completed!"
        break
    elif [ "$status" = "failed" ]; then
        echo "❌ Test failed!"
        break
    fi
    
    sleep 5
    attempt=$((attempt + 1))
done

# 7. Get final results
echo "7. Getting final results..."
response=$(curl -s -w "%{http_code}" "$SERVER_URL/test/$test_id")
status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

# 8. Display summary if test completed
if echo "$body" | jq -e '.overall' >/dev/null 2>&1; then
    echo "📊 Test Summary:"
    echo "   Total Requests: $(echo "$body" | jq -r '.overall.total_requests')"
    echo "   Total Success: $(echo "$body" | jq -r '.overall.total_success')"
    echo "   Overall RPS: $(echo "$body" | jq -r '.overall.overall_rps')"
    echo "   Success Rate: $(echo "$body" | jq -r '.overall.overall_success_rate')%"
    
    # Display method-specific results
    echo "🔍 Method-Specific Results:"
    echo "$body" | jq -r '.results[]? | "   📈 \(.method_name): Requests: \(.total_requests), Success Rate: \(.success_rate)%, RPS: \(.requests_per_sec)"' 2>/dev/null
fi

# 9. Delete the tests
echo "9. Deleting tests..."
response=$(curl -s -w "%{http_code}" \
  -X DELETE \
  "$SERVER_URL/test/$test_id")

status_code="${response: -3}"
body="${response%???}"
print_response "$status_code" "$body"

if [ "$test_id2" != "null" ] && [ -n "$test_id2" ]; then
    response=$(curl -s -w "%{http_code}" \
      -X DELETE \
      "$SERVER_URL/test/$test_id2")
    
    status_code="${response: -3}"
    body="${response%???}"
    print_response "$status_code" "$body"
fi

echo "✅ Examples completed!" 