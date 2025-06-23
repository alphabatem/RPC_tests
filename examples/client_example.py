#!/usr/bin/env python3
"""
RPC Test Server Client Example

This script demonstrates how to use the RPC Test Server API to run tests
and retrieve results with method-specific configurations.
"""

import requests
import json
import time
import sys

# Server configuration
SERVER_URL = "http://localhost:8080"

def print_response(response, title="Response"):
    """Pretty print API response"""
    print(f"\n{'='*50}")
    print(f"{title}")
    print(f"{'='*50}")
    print(f"Status Code: {response.status_code}")
    print(f"Headers: {dict(response.headers)}")
    print(f"Body:")
    print(json.dumps(response.json(), indent=2, default=str))

def start_test():
    """Start a new RPC test with method-specific configurations"""
    print("🚀 Starting a new RPC test with method-specific configurations...")
    
    # Test configuration with method-specific settings
    test_config = {
        "remote_rpc_url": "https://us.rpc.fluxbeam.xyz",
        "rpc_apikey": "YOUR_API_KEY_HERE",  # Replace with your actual API key
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
                "enabled": True
            },
            "getMultipleAccounts": {
                "concurrency": 5,
                "duration": 15,
                "limit": 100,
                "enabled": True
            },
            "getProgramAccounts": {
                "concurrency": 3,
                "duration": 10,
                "limit": 25,
                "enabled": False  # Disable this method
            }
        }
    }
    
    try:
        response = requests.post(
            f"{SERVER_URL}/test",
            json=test_config,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            result = response.json()
            test_id = result.get("test_id")
            print(f"✅ Test started successfully!")
            print(f"📋 Test ID: {test_id}")
            return test_id
        else:
            print(f"❌ Failed to start test: {response.status_code}")
            print_response(response, "Error Response")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Network error: {e}")
        return None

def start_simple_test():
    """Start a simple test using global defaults"""
    print("🚀 Starting a simple RPC test with global defaults...")
    
    # Simple test configuration using global defaults
    test_config = {
        "target_rpc_url": "https://api.mainnet-beta.solana.com",
        "global_config": {
            "concurrency": 3,
            "duration": 10
        }
        # Methods will use global config defaults
    }
    
    try:
        response = requests.post(
            f"{SERVER_URL}/test",
            json=test_config,
            headers={"Content-Type": "application/json"}
        )
        
        if response.status_code == 200:
            result = response.json()
            test_id = result.get("test_id")
            print(f"✅ Simple test started successfully!")
            print(f"📋 Test ID: {test_id}")
            return test_id
        else:
            print(f"❌ Failed to start simple test: {response.status_code}")
            print_response(response, "Error Response")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Network error: {e}")
        return None

def get_test_status(test_id):
    """Get the status of a running test"""
    try:
        response = requests.get(f"{SERVER_URL}/test/{test_id}")
        
        if response.status_code == 200:
            result = response.json()
            status = result.get("status", "unknown")
            print(f"📊 Test {test_id} status: {status}")
            return status
        else:
            print(f"❌ Failed to get test status: {response.status_code}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Network error: {e}")
        return None

def get_test_results(test_id):
    """Get the results of a completed test"""
    try:
        response = requests.get(f"{SERVER_URL}/test/{test_id}")
        
        if response.status_code == 200:
            result = response.json()
            if result.get("status") == "completed":
                print_response(response, f"Test Results for {test_id}")
                return result
            else:
                print(f"⏳ Test {test_id} is still running...")
                return None
        else:
            print(f"❌ Failed to get test results: {response.status_code}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Network error: {e}")
        return None

def list_all_tests():
    """List all tests"""
    try:
        response = requests.get(f"{SERVER_URL}/tests")
        
        if response.status_code == 200:
            print_response(response, "All Tests")
            return response.json()
        else:
            print(f"❌ Failed to list tests: {response.status_code}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Network error: {e}")
        return None

def delete_test(test_id):
    """Delete a test"""
    try:
        response = requests.delete(f"{SERVER_URL}/test/{test_id}")
        
        if response.status_code == 200:
            print(f"✅ Test {test_id} deleted successfully!")
            return True
        else:
            print(f"❌ Failed to delete test: {response.status_code}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Network error: {e}")
        return False

def wait_for_test_completion(test_id, max_wait_time=300):
    """Wait for a test to complete"""
    print(f"⏳ Waiting for test {test_id} to complete...")
    start_time = time.time()
    
    while time.time() - start_time < max_wait_time:
        status = get_test_status(test_id)
        
        if status == "completed":
            print(f"✅ Test {test_id} completed!")
            return True
        elif status == "failed":
            print(f"❌ Test {test_id} failed!")
            return False
        elif status is None:
            print(f"❌ Could not get status for test {test_id}")
            return False
        
        # Wait 5 seconds before checking again
        time.sleep(5)
    
    print(f"⏰ Timeout waiting for test {test_id} to complete")
    return False

def main():
    """Main function demonstrating the API usage"""
    print("🔧 RPC Test Server Client Example")
    print("="*50)
    
    # Check if server is running
    try:
        response = requests.get(f"{SERVER_URL}/")
        if response.status_code != 200:
            print(f"❌ Server is not responding properly: {response.status_code}")
            sys.exit(1)
        print("✅ Server is running!")
    except requests.exceptions.RequestException as e:
        print(f"❌ Cannot connect to server at {SERVER_URL}")
        print(f"   Make sure the server is running with: ./rpc_test server")
        sys.exit(1)
    
    # List existing tests
    print("\n📋 Listing existing tests...")
    list_all_tests()
    
    # Start a test with method-specific configurations
    test_id = start_test()
    if not test_id:
        print("❌ Failed to start test, exiting...")
        sys.exit(1)
    
    # Wait for test completion
    if wait_for_test_completion(test_id):
        # Get results
        results = get_test_results(test_id)
        if results:
            print("\n🎉 Test completed successfully!")
            
            # Display summary
            overall = results.get("overall", {})
            if overall:
                print(f"\n📊 Overall Results:")
                print(f"   Total Requests: {overall.get('total_requests', 0)}")
                print(f"   Total Success: {overall.get('total_success', 0)}")
                print(f"   Overall RPS: {overall.get('overall_rps', 0):.2f}")
                print(f"   Success Rate: {overall.get('overall_success_rate', 0):.2f}%")
            
            # Display method-specific results
            method_results = results.get("results", [])
            if method_results:
                print(f"\n🔍 Method-Specific Results:")
                for result in method_results:
                    method_name = result.get("method_name", "unknown")
                    print(f"   📈 {method_name}:")
                    print(f"      Requests: {result.get('total_requests', 0)}")
                    print(f"      Success Rate: {result.get('success_rate', 0):.2f}%")
                    print(f"      RPS: {result.get('requests_per_sec', 0):.2f}")
    
    # List tests again to see the completed test
    print("\n📋 Listing tests after completion...")
    list_all_tests()
    
    # Clean up - delete the test
    print(f"\n🧹 Cleaning up test {test_id}...")
    delete_test(test_id)
    
    print("\n✅ Example completed successfully!")

if __name__ == "__main__":
    main() 