package service_test

import (
	"bitbucket.org/wepala/weos-controller/service"
	"bytes"
	"flag"
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

func TestServeHTTP(t *testing.T){
	var handler http.Handler
	//setup html server

	controllerService, _ := service.NewControllerService("testdata/api/x-mock-status-code.yaml", service.NewPluginLoader())
	handler = service.NewHTTPServer(controllerService, "static")

	rw := httptest.NewRecorder()

	//send test request
	log.Debugf("Load input fixture: %s", "x_mock_status_code.input.http")
	request := loadHttpRequestFixture(filepath.Join("testdata/html/http", "x_mock_status_code.input.http"), t)
	handler.ServeHTTP(rw, request)

	response := rw.Result()
	statusCode := strconv.Itoa(response.StatusCode)

	if response == nil{
		t.Error("Response expected but returned nothing")
	}

	if request.Header.Get("X-MOCK-STATUS-CODE") != statusCode{
		t.Errorf("Expected response code %s, got %s instead", request.Header.Get("X-MOCK-STATUS-CODE"), statusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	newStr := buf.Bytes()

	if string(newStr) != "<html>\n  <head>\n      <title>X Mock Example Page</title>\n  </head>\n  <body>\n    This is a mocked page\n  </body>\n</html>\n"{
		t.Errorf("Incorrect example returned, expected %s, but got %s", "<html>\n  <head>\n      <title>X Mock Example Page</title>\n  </head>\n  <body>\n    This is a mocked page\n  </body>\n</html>\n", newStr)
	}
}