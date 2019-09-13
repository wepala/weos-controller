/**
 * Test Helpers
 */
package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"github.com/wepala/go-testhelpers"
	"io/ioutil"
	"net/http"
	"testing"
)

type HTTPTest struct {
	name          string
	testDataDir   string
	apiFixture    string
	configFixture string
	service       service.ServiceInterface
}

//loadHttpRequestFixture wrapper around the test helper to make it easier to use it with test table
func loadHttpRequestFixture(filename string, t *testing.T) *http.Request {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}

	request, err := testhelpers.LoadHTTPRequestFixture(data)
	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}

	return request
}

//loadHttpResponseFixture wrapper around the test helper to make it easier to use it with test table
func loadHttpResponseFixture(filename string, request *http.Request, t *testing.T) *http.Response {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("text fixture not loaded %v", err)
	}

	response, err := testhelpers.LoadHTTPResponseFixture(data, request)
	if err != nil {
		t.Fatalf("text fixture not loaded %v", err)
	}

	return response
}
