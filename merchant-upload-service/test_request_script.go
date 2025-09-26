//go:build ignore
// +build ignore

// This is a standalone test script
// To run: go run test_request_script.go

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// This is a standalone script and should be run separately from the main application
// It is not part of the main application codebase
func main() {
	// Use the exact URL that's causing issues
	url := "http://localhost:4300/merchant/1000/file?path=www%252Fmerchant_1000%252Fdata%252Fdomains.json"
	fmt.Printf("Making request to: %s\n", url)
	
	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	
	// Print the status code and response body
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", body)
}
