package service_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"bitbucket.org/wepala/weos-controller/service"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

//go:generate moq -out testing_mocks_test.go -pkg service_test . ServiceInterface
/*
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
}*/

var httpServerTests = []*HTTPTest{
	{
		name:        "about_page_200",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/basic-site-api.yml",
	},
	{
		name:        "x_mock_status_code",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/x-mock-status-code.yaml",
	},
	{
		name:        "x_mock_multiple_examples",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/x-mock-status-code.yaml",
	},
	{
		name:        "x_mock_component_example",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/x-mock-status-code.yaml",
	},
	{
		name:        "x_mock_array_component_example",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/x-mock-status-code.yaml",
	},
	{
		name:        "wildcard",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/http-test-api-wildcard.yaml",
	},
	{
		name:        "wildcard_nested",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/http-test-api-wildcard.yaml",
	},
	{
		name:        "wildcard_root",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/http-test-api-wildcard.yaml",
	},
	{
		name:        "wildcard_static",
		testDataDir: "testdata/html/http",
		apiFixture:  "testdata/api/http-test-api-wildcard-static.yaml",
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
	t.SkipNow()
	runHttpServerTests(httpServerTests, false, "static", t)
}

func runHttpServerTests(tests []*HTTPTest, serveStatic bool, staticFolder string, t *testing.T) {
	//t.SkipNow()
	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			var handler http.Handler
			//setup html server
			controllerService, _ := service.NewControllerService(test.apiFixture, service.NewPluginLoader())
			handler = service.NewHTTPServer(controllerService, serveStatic, staticFolder)

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
	//t.SkipNow()
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

		if database.Id != "35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "someid", database.Id)
		}

		if database.Wern != "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "somewern", database.Wern)
		}
	})

	t.Run("test return specific content type", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_content_type.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_content_type.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if strconv.Itoa(rw.Result().StatusCode) != "200" {
			t.Errorf("expected the response code to be %s, got %d", "200", rw.Result().StatusCode) //changed as there was no status code defined in input file
		}

		if rw.Result().Header.Get("Content-Type") != "text/html" {
			t.Errorf("expected the Content-Type to be %s, got %s", "text/html", rw.Result().Header.Get("Content-Type"))
		}

		//confirm the body
		expectedBody := "test"
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedBody) {
			t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(expectedBody), strings.TrimSpace(string(body)))
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

		if database[0].Id != "35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "35a54035-753d-4123-bea2-ff3ee25b0eea", database[0].Id)
		}

		if database[0].Wern != "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea", database[0].Wern)
		}
	})

	t.Run("test that CORs headers are NOT automatically set", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		if rw.Result().Header.Get("Access-Control-Allow-Origin") != "" {
			t.Error("no response headers was expected")
		}
	})

	t.Run("test defined headers are returned in response", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_header_example.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_header_example.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/databases"),
		}

		mockHandler.ServeHTTP(rw, request)

		if rw.Result().Header.Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("expected header Access-Control-Allow-Origin to be %s, got %s", "*", rw.Result().Header.Get("Access-Control-Allow-Origin"))
		}

		if rw.Result().Header.Get("Access-Control-Allow-Headers") != "Authorization, Content-Type" {
			t.Errorf("expected header Access-Control-Allow-Headers to be %s, got %s", "Authorization, Content-Type", rw.Result().Header.Get("Access-Control-Allow-Headers"))
		}
	})

	t.Run("test no status code hits default", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/databases"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if rw.Result().StatusCode != 200 {
			t.Errorf("expected the response code to be %d, got %d", 200, rw.Result().StatusCode)
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

		if database.Id != "35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "default", database.Id)
		}

		if database.Wern != "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "default", database.Wern)
		}
	})

	t.Run("test no status code hits 200", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_no_status_code_no_default.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code_no_default.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/nodefault"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if rw.Result().StatusCode != 200 {
			t.Errorf("expected the response code to be %d, got %d", 200, rw.Result().StatusCode)
		}

		database := &struct {
			Id   string `json:"id"`
			Wern string `json:"wern"`
		}{}

		err := json.Unmarshal(body, database)
		if err != nil {
			t.Errorf("expected json response, %q", err.Error())
		}

		if database.Id != "35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "default", database.Id)
		}

		if database.Wern != "weos:tt:data:12345:35a54035-753d-4123-bea2-ff3ee25b0eea" {
			t.Errorf("expected the id on the response to be %s, got %s", "default", database.Wern)
		}
	})

}

func TestMockHandler_ServeHTTPErrors(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	config, err := loader.LoadSwaggerFromFile("testdata/api/x-mock-status-code.yaml")

	if err != nil {
		t.Fatalf("error loading %s: %s", "testdata/api/x-mock-status-code.yaml", err.Error())
	}

	t.Run("test undefined status code ", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_missing_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_missing_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if strconv.Itoa(rw.Result().StatusCode) == request.Header.Get("X-Mock-Status-Code") {
			t.Errorf("expected the response code to be %s, got %d", "200", rw.Result().StatusCode)
		}

		if rw.Result().Header.Get("Content-Type") != "text/plain" {
			t.Errorf("expected the Content-Type to be %s, got %s", "text/plain", rw.Result().Header.Get("Content-Type"))
		}

		//confirm the body
		expectedBody := fmt.Sprintf("There is no mocked response for status code %s", request.Header.Get("X-Mock-Status-Code"))
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedBody) {
			t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
		}

	})

	t.Run("test undefined example", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_multiple_examples_missing_example.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_multiple_examples_missing_example.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/about"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		if rw.Result().Header.Get("Content-Type") != "text/html" {
			t.Errorf("expected the Content-Type to be %s, got %s", "text/html", rw.Result().Header.Get("Content-Type"))
		}

		//confirm the body
		expectedBody := fmt.Sprintf("There is no mocked response with example named %s", request.Header.Get("X-Mock-Example"))
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedBody) {
			t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
		}

	})

	t.Run("test multiple examples none specified", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_no_example.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_example.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/about"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		//confirm the body
		expectedBody := "There are multiple examples defined. Please specify one using the X-Mock-Example header"
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedBody) {
			t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
		}

	})

	t.Run("test multiple content types defined none specified", func(t *testing.T) {
		log.Debugf("Load input fixture: %s", "x_mock_no_content_type.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_content_type.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		//confirm the body
		expectedBody := "There are multiple content types defined. Please specify one using the X-Mock-Content-Type header"
		if strings.TrimSpace(string(body)) != strings.TrimSpace(expectedBody) {
			t.Errorf("expected body '%s', got: '%s'", strings.TrimSpace(string(expectedBody)), strings.TrimSpace(string(body)))
		}

	})

}

func TestOtherSwaggerFiles(t *testing.T) {
	t.Run("test callback example", func(t *testing.T) {
		loader := openapi3.NewSwaggerLoader()
		config, err := loader.LoadSwaggerFromFile("testdata/api/callback-example.yaml")

		if err != nil {
			t.Fatalf("error loading %s: %s", "testdata/api/callback-example.yaml", err.Error())
		}

		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/streams"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		log.Info(string(body))

		if string(body) == "" {
			t.Error("Expected something to be returned")
		}
	})

	t.Run("test api example", func(t *testing.T) {
		loader := openapi3.NewSwaggerLoader()
		config, err := loader.LoadSwaggerFromFile("testdata/api/api-with-examples.yaml")

		if err != nil {
			t.Fatalf("error loading %s: %s", "testdata/api/api-with-examples.yaml", err.Error())
		}

		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		log.Info(string(body))

		if string(body) == "" {
			t.Error("Expected something to be returned")
		}
	})

	t.Run("test petstore expanded example", func(t *testing.T) {
		loader := openapi3.NewSwaggerLoader()
		config, err := loader.LoadSwaggerFromFile("testdata/api/petstore-expanded.yaml")

		if err != nil {
			t.Fatalf("error loading %s: %s", "testdata/api/petstore-expanded.yaml", err.Error())
		}

		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/pets"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		log.Info(string(body))

		if string(body) == "" {
			t.Error("Expected something to be returned")
		}
	})

	t.Run("test petstore example", func(t *testing.T) {
		loader := openapi3.NewSwaggerLoader()
		config, err := loader.LoadSwaggerFromFile("testdata/api/petstore.yaml")

		if err != nil {
			t.Fatalf("error loading %s: %s", "testdata/api/petstore.yaml", err.Error())
		}

		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/pets"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		log.Info(string(body))

		if string(body) == "" {
			t.Error("Expected something to be returned")
		}
	})

	t.Run("test uspto example", func(t *testing.T) {
		loader := openapi3.NewSwaggerLoader()
		config, err := loader.LoadSwaggerFromFile("testdata/api/uspto.yaml")

		if err != nil {
			t.Fatalf("error loading %s: %s", "testdata/api/uspto.yaml", err.Error())
		}

		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		log.Info(string(body))

		if string(body) == "" {
			t.Error("Expected something to be returned")
		}
	})

	t.Run("test link example", func(t *testing.T) {
		loader := openapi3.NewSwaggerLoader()
		config, err := loader.LoadSwaggerFromFile("testdata/api/link-example.yaml")

		if err != nil {
			t.Fatalf("error loading %s: %s", "testdata/api/link-example.yaml", err.Error())
		}

		log.Debugf("Load input fixture: %s", "x_mock_no_status_code.input.http")
		request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_no_status_code.input.http"), t)
		rw := httptest.NewRecorder()

		mockHandler := service.MockHandler{
			PathInfo: config.Paths.Find("/2.0/users/{username}"),
		}

		mockHandler.ServeHTTP(rw, request)

		body, _ := ioutil.ReadAll(rw.Result().Body)

		log.Info(string(body))

		if string(body) == "" {
			t.Error("Expected something to be returned")
		}
	})
}

func Test_WECON_2(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	config, err := loader.LoadSwaggerFromFile("testdata/api/x-mock-status-code.yaml")

	if err != nil {
		t.Fatalf("error loading mock api config %s", err)
	}

	request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "about_page_200.input.http"), t)
	t.Run("first and second handler called, second handler responds", func(t *testing.T) {
		rw := httptest.NewRecorder()
		mockHandler1 := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		})

		mockHandler2 := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(200)
			rw.Write([]byte(`bar`))
		})

		mockService := &ServiceInterfaceMock{
			GetConfigFunc: func() *openapi3.Swagger {
				return config
			},
			GetGlobalMiddlewareConfigFunc: func() (configs []*service.MiddlewareConfig, err error) {
				return nil, nil
			},
			GetHandlersFunc: func(config *service.PathConfig, mockHandler http.Handler) (funcs []http.HandlerFunc, err error) {
				return []http.HandlerFunc{mockHandler1, mockHandler2}, nil
			},
			GetPathConfigFunc: func(path string, operation string) (config *service.PathConfig, err error) {
				return nil, nil
			},
		}

		httpapi := service.NewHTTPServer(mockService, false, "")
		httpapi.ServeHTTP(rw, request)
		//check to see if the first status code is registered
		response := rw.Result()
		if response.StatusCode != 200 {
			t.Errorf("expected the status code to be %d, got %d", 200, response.StatusCode)
		}

		result, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("error retreiving response %s", err)
		}

		if string(result) != "bar" {
			t.Errorf("expected result to be %s, got %s", "bar", string(result))
		}
	})
	t.Run("second handler no response if first handler responses", func(t *testing.T) {
		rw := httptest.NewRecorder()
		mockHandler1 := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(500)
			rw.Write([]byte(`foo`))
		})

		mockHandler2 := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(200)
			rw.Write([]byte(`bar`))
		})

		mockService := &ServiceInterfaceMock{
			GetConfigFunc: func() *openapi3.Swagger {
				return config
			},
			GetGlobalMiddlewareConfigFunc: func() (configs []*service.MiddlewareConfig, err error) {
				return nil, nil
			},
			GetHandlersFunc: func(config *service.PathConfig, mockHandler http.Handler) (funcs []http.HandlerFunc, err error) {
				return []http.HandlerFunc{mockHandler1, mockHandler2}, nil
			},
			GetPathConfigFunc: func(path string, operation string) (config *service.PathConfig, err error) {
				return nil, nil
			},
		}

		httpapi := service.NewHTTPServer(mockService, false, "")
		httpapi.ServeHTTP(rw, request)
		//check to see if the first status code is registered
		response := rw.Result()
		if response.StatusCode != 500 {
			t.Errorf("expected the status code to be %d, got %d", 500, response.StatusCode)
		}

		result, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("error retreiving response %s", err)
		}

		if string(result) != "foo" {
			t.Errorf("expected result to be %s, got %s", "foo", string(result))
		}
	})

}
