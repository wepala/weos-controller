package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"encoding/json"
	"flag"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface

var mockServerTests = []*HTTPTest{
	{
		name:        "about_page_OPTION",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api.yml",
	},
	{
		name:        "landingpage_mock_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api.yml",
	},
	{
		name:        "poll_list_mock_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/rest-api.yml",
	},
	{
		name:        "apollo_list_mock_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/apollo-api.yaml",
	},
}

var httpServerTests = []*HTTPTest{
	{
		name:        "about_page_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api." + runtime.GOOS + ".yml",
	},
}

var staticPageTest = []*HTTPTest{
	{
		name:        "about_page_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api.yml",
	},
}

var update = flag.Bool("update", false, "update .golden files")

func Test_Endpoints(t *testing.T) {
	runHttpServerTests(httpServerTests, "static", t)
}

func runHttpServerTests(tests []*HTTPTest, staticFolder string, t *testing.T) {
	//t.SkipNow()
	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			var handler http.Handler
			//setup html server
			controllerService, _ := service.NewControllerService(test.apiFixture, service.NewPluginLoader())
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

func TestMockHandler_ServeHTTP(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	config, err := loader.LoadSwaggerFromFile("testdata/api/x-mock-status-code.yaml")

	if err != nil {
		t.Fatalf("error loading %s: %s", "testdata/api/x-mock-status-code.yaml", err.Error())
	}

	t.Run("test basic example", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if strconv.Itoa(rw.Result().StatusCode) != request.Header.Get("X-Mock-Status-Code") {
			t.Errorf("expected the response code to be %s, got %d", request.Header.Get("X-Mock-Status-Code"), rw.Result().StatusCode)
		}

		if rw.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected the Content-Type to be %s, got %s", "application/json", rw.Result().Header.Get("Content-Type"))
		}

		database := &struct {
			Id   string `json:"id"`
			Wern string `json:"wern"`
		}{}

		err := json.Unmarshal(body, database)
		if err != nil {
			t.Errorf("expected json response, %q", err.Error())
		}

		if database.Id != "someid" {
			t.Errorf("expected the id on the response to be %s, got %s", "someid", database.Id)
		}

		if database.Wern != "somewern" {
			t.Errorf("expected the id on the response to be %s, got %s", "somewern", database.Wern)
		}
	})

	t.Run("test multiple examples", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_multiple_examples.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_multiple_examples.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/about"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)
		expectedResponse := loadHttpResponseFixture(filepath.Join("testdata/html/http", "x_mock_multiple_examples.golden.http"), request, t)

		if strconv.Itoa(rw.Result().StatusCode) != request.Header.Get("X-Mock-Status-Code") {
			t.Errorf("expected the response code to be %s, got %d", request.Header.Get("X-Mock-Status-Code"), rw.Result().StatusCode)
		}

		if rw.Result().Header.Get("Content-Type") != "text/html" {
			t.Errorf("expected the Content-Type to be %s, got %s", "text/html", rw.Result().Header.Get("Content-Type"))
		}

		//confirm the body
		expectedBody, _ := ioutil.ReadAll(expectedResponse.Body)
		if strings.TrimSpace(string(body)) != strings.TrimSpace(string(expectedBody)) {
			t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
		}
	})

	t.Run("test example on component", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_component_example.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_component_example.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/databases"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if strconv.Itoa(rw.Result().StatusCode) != request.Header.Get("X-Mock-Status-Code") {
			t.Errorf("expected the response code to be %s, got %d", request.Header.Get("X-Mock-Status-Code"), rw.Result().StatusCode)
		}

		if rw.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected the Content-Type to be %s, got %s", "application/json", rw.Result().Header.Get("Content-Type"))
		}

		database := &struct {
			Id   string `json:"id"`
			Wern string `json:"wern"`
		}{}

		if rw.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected the Content-Type to be %s, got %s", "application/json", rw.Result().Header.Get("Content-Type"))
		}

		err := json.Unmarshal(body, database)
		if err != nil {
			t.Errorf("expected json response, %q", err.Error())
		}

		if database.Id != "35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "35a54035-753d-4123-bea2-ff3ee25b0eea", database.Id)
		}

		if database.Wern != "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea", database.Wern)
		}
	})

	t.Run("test example on component when response is an array", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_array_component_example.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_array_component_example.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/databases"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if strconv.Itoa(rw.Result().StatusCode) != request.Header.Get("X-Mock-Status-Code") {
			t.Errorf("expected the response code to be %s, got %d", request.Header.Get("X-Mock-Status-Code"), rw.Result().StatusCode)
		}

		if rw.Result().Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected the Content-Type to be %s, got %s", "application/json", rw.Result().Header.Get("Content-Type"))
		}

		var database []*struct {
			Id   string `json:"id"`
			Wern string `json:"wern"`
		}

		err := json.Unmarshal(body, &database)
		if err != nil {
			t.Errorf("expected json response, %q", err.Error())
		}

		if strconv.Itoa(len(database)) != request.Header.Get("X-Mock-Example-Length") {
			t.Errorf("expected the length of the result to be %s, got %d", request.Header.Get("X-Mock-Example-Length"), len(database))
		}

		if database[5].Id != "35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "35a54035-753d-4123-bea2-ff3ee25b0eea", database[0].Id)
		}

		if database[5].Wern != "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea", database[0].Wern)
		}
	})

}
