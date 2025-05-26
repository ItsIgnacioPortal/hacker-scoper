package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
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
