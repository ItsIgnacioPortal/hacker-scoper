package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/zenizh/go-capturer"
)

//========================================================================
//                            HELPER FUNCTIONS
//========================================================================

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func checkForErrors(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

//========================================================================
//========================================================================
//========================================================================

func Test_parseOutOfScopes(t *testing.T) {
	// Simple test - inscope URL
	assetURL, _ := url.Parse("https://example.com")
	outOfScopeString := "zendesk*.example.com"
	value := parseOutOfScopes(assetURL, outOfScopeString, nil)
	equals(t, false, value)

	// Simple test - out of scope URL
	assetURL, _ = url.Parse("https://zendesk.internal.example.com")
	outOfScopeString = "zendesk*.example.com"
	value = parseOutOfScopes(assetURL, outOfScopeString, nil)
	equals(t, true, value)

	// Simple test - in-scope URL with a URL-like out-of-scope string
	assetURL, _ = url.Parse("https://zendesk.internal.example.com")
	outOfScopeString = "https://sometool.internal.example.com"
	value = parseOutOfScopes(assetURL, outOfScopeString, nil)
	equals(t, false, value)

	// Simple test - out-of-scope URL with a URL-like out-of-scope string
	assetURL, _ = url.Parse("https://zendesk.internal.example.com")
	outOfScopeString = "https://zendesk.internal.example.com"
	value = parseOutOfScopes(assetURL, outOfScopeString, nil)
	equals(t, true, value)

	// Test - in-scope URL with a URL-like out-of-scope string with an unusual scheme
	assetURL, _ = url.Parse("https://zendesk.internal.example.com")
	outOfScopeString = "mongodb://sometool.internal.example.com"
	value = parseOutOfScopes(assetURL, outOfScopeString, nil)
	equals(t, false, value)

	// Test - out-of-scope URL with a URL-like out-of-scope string with an unusual scheme
	assetURL, _ = url.Parse("https://zendesk.internal.example.com")
	outOfScopeString = "mongodb://zendesk.internal.example.com"
	value = parseOutOfScopes(assetURL, outOfScopeString, nil)
	equals(t, true, value)

	// Test with a bad function invocation, providing both an assetURL and an assetIP
	// Only the assetURL should be used in this case
	assetURL, _ = url.Parse("https://zendesk.internal.example.com")
	outOfScopeString = "zendesk*.example.com"
	assetIP := net.ParseIP("127.0.0.1")
	value = parseOutOfScopes(assetURL, outOfScopeString, assetIP)
	equals(t, true, value)

	// Test with a bad function invocation, providing both assetURL and assetIP as nil
	// The function should return false
	assetURL = nil
	outOfScopeString = "127.0.0.1"
	assetIP = nil
	value = parseOutOfScopes(assetURL, outOfScopeString, assetIP)
	equals(t, false, value)
}

func Example_parseOutOfScopes() {
	// Test with an invalid out-of-scope string
	// In context, this function would print a warning to stderr and return false
	// However, for testing purposes, we will just check the stederr output
	assetURL, _ := url.Parse("https://example.com")
	outOfScopeString := "this is not even close to a URL"

	out := capturer.CaptureStderr(func() {
		_ = parseOutOfScopes(assetURL, outOfScopeString, nil)
	})

	fmt.Println(out)
	// Output: [33m[WARNING]: Couldn't parse out-of-scope "[38;2;0;204;255mhttps://[33mthis is not even close to a URL" as a URL.[0m
}

func Test_updateFireBountyJSON(t *testing.T) {
	// This test just verifies if the firebountyAPIURL is still available online, and if the JSON it returns still matches the expected structure.
	// firebountyAPIURL is a global variable defined in the main package.
	// First, we test if the URL is reachable with a HEAD request.
	fmt.Println(firebountyAPIURL)
	resp, err := http.Head("https://firebounty.com/api/v1/scope/all/url_only/")
	// if error is not nil and the response body has more than 1 byte, we fail the test.
	if err != nil || resp == nil || resp.ContentLength < 1 {
		t.Fatalf("Failed to reach firebounty API URL: %v", err)
	} else {
		// If the HEAD request is successful, we proceed to test the JSON structure.
		// We can use a simple HTTP GET request to fetch the JSON.
		resp, err = http.Get(firebountyAPIURL)
		checkForErrors(t, err)
		defer resp.Body.Close()

		// We can check if the Content-Type is application/json
		if resp.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Expected Content-Type application/json, got %s", resp.Header.Get("Content-Type"))
		}

		// We can also check if the response body is not empty
		if resp.ContentLength == 0 {
			t.Fatal("Expected non-empty response body")
		}
	}
}

func Test_removePortFromHost(t *testing.T) {
	// testURL must be in a variable of type *url.URL, which contains "https://example.com:8080/path?query=123"
	testURL, _ := url.Parse("https://example.com:8080/path?query=123")
	value := removePortFromHost(testURL)
	equals(t, "example.com", value)
}

func Test_removeDuplicateStr(t *testing.T) {
	// testSlice must be a slice of strings with duplicates
	testSlice := []string{"a", "b", "a", "c", "b"}
	value := removeDuplicateStr(testSlice)
	equals(t, []string{"a", "b", "c"}, value)
}
