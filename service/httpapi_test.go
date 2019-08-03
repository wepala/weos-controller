package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"flag"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"path/filepath"
	"strings"
	"testing"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface

var mockServerTests = []*HTTPTest{
	{
		name:        "landingpage_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api.yml",
	},
	{
		name:        "poll_list_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/rest-api.yml",
	},
}

var httpServerTests = []*HTTPTest{
	{
		name:        "landingpage_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api.yml",
	},
	{
		name:        "poll_list_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/rest-api.yml",
	},
}

func Test_Endpoints(t *testing.T) {
	runMockServerTests(mockServerTests, "static", t)
}

func runMockServerTests(tests []*HTTPTest, staticFolder string, t *testing.T) {
	var update = flag.Bool("update", false, "update .golden files")
	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			var handler http.Handler
			//setup html server
			controllerService, _ := service.NewControllerService(test.apiFixture, "", nil)
			handler = service.NewMockHTTPServer(controllerService, staticFolder)

			rw := httptest.NewRecorder()

			//send test request
			log.Debugf("Load input fixture: %s", test.name+".input.http")
			request := loadHttpRequestFixture(filepath.Join(test.testDataDir, test.name+".input.http"), t)
			handler.ServeHTTP(rw, request)

			//confirm response
			response := rw.Result()

			responseFixture := filepath.Join(test.testDataDir, test.name+".golden.http")
			if *update {
				responseBytes, _ := httputil.DumpResponse(response, true)
				err := ioutil.WriteFile(responseFixture, responseBytes, 0644)
				if err != nil {
					t.Fatalf("failed to write fixture '%s' with error '%v'", responseFixture, err)
				}
			}
			body, _ := ioutil.ReadAll(response.Body)
			expectedResponse := loadHttpResponseFixture(responseFixture, request, t)

			//confirm the expected status code
			if response.StatusCode != expectedResponse.StatusCode {
				t.Errorf("expected status code %d, got: %d", expectedResponse.StatusCode, response.StatusCode)
			}

			//confirm the content type returned
			if response.Header.Get("Content-Type") != expectedResponse.Header.Get("Content-Type") {
				t.Errorf("expected content type %s, got: %s", expectedResponse.Header.Get("Content-Type"), response.Header.Get("Content-Type"))
			}

			//confirm the body
			expectedBody, _ := ioutil.ReadAll(expectedResponse.Body)
			if strings.TrimSpace(string(body)) != strings.TrimSpace(string(expectedBody)) {
				t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
			}
		})
	}
}

func runHttpServerTests(tests []*HTTPTest, staticFolder string, t *testing.T) {
	var update = flag.Bool("update", false, "update .golden files")
	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			var handler http.Handler
			//setup html server
			controllerService, _ := service.NewControllerService(test.apiFixture, "", nil)
			handler = service.NewHTTPServer(controllerService, staticFolder)

			rw := httptest.NewRecorder()

			//send test request
			log.Debugf("Load input fixture: %s", test.name+".input.http")
			request := loadHttpRequestFixture(filepath.Join(test.testDataDir, test.name+".input.http"), t)
			handler.ServeHTTP(rw, request)

			//confirm response
			response := rw.Result()

			responseFixture := filepath.Join(test.testDataDir, test.name+".golden.http")
			if *update {
				responseBytes, _ := httputil.DumpResponse(response, true)
				err := ioutil.WriteFile(responseFixture, responseBytes, 0644)
				if err != nil {
					t.Fatalf("failed to write fixture '%s' with error '%v'", responseFixture, err)
				}
			}
			body, _ := ioutil.ReadAll(response.Body)
			expectedResponse := loadHttpResponseFixture(responseFixture, request, t)

			//confirm the expected status code
			if response.StatusCode != expectedResponse.StatusCode {
				t.Errorf("expected status code %d, got: %d", expectedResponse.StatusCode, response.StatusCode)
			}

			//confirm the content type returned
			if response.Header.Get("Content-Type") != expectedResponse.Header.Get("Content-Type") {
				t.Errorf("expected content type %s, got: %s", expectedResponse.Header.Get("Content-Type"), response.Header.Get("Content-Type"))
			}

			//confirm the body
			expectedBody, _ := ioutil.ReadAll(expectedResponse.Body)
			if strings.TrimSpace(string(body)) != strings.TrimSpace(string(expectedBody)) {
				t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
			}
		})
	}
}
