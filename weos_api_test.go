package weoscontroller_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	weoscontroller "github.com/wepala/weos-controller"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAPI_RequestRecording(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/endpoint", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	api := &weoscontroller.API{Config: &weoscontroller.APIConfig{
		RecordingBaseFolder: ".",
	}}
	e.PUT("/endpoint", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}, api.RequestRecording)
	api.SetEchoInstance(e)
	e.ServeHTTP(rec, req)
	//confirm file is created
	if _, err := os.Stat("./_endpoint.input.http"); os.IsNotExist(err) {
		t.Fatalf("expected fixture file to be created")
	}

	//confirm contents of file
	request := loadHttpRequestFixture("./_endpoint.input.http", t)
	var bodyStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(request.Body).Decode(&bodyStruct)
	if err != nil {
		t.Fatalf("error unmarshalling request '%s'", err)
	}

	if bodyStruct.Name != "Sojourner Truth" {
		t.Errorf("expected the name to be '%s', got '%s'", "Sojourner Truth", bodyStruct.Name)
	}

	//delete file
	err = os.Remove("./_endpoint.input.http")
	if err != nil {
		t.Fatalf("unable to delete test fixture created '%s', got error '%s'", "./_endpoint.input.http", err)
	}
}

//loadHttpRequestFixture wrapper around the test helper to make it easier to use it with test table
func loadHttpRequestFixture(filename string, t *testing.T) *http.Request {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}

	reader := bufio.NewReader(bytes.NewReader(data))
	request, err := http.ReadRequest(reader)
	if err == io.EOF {
		return request
	}

	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}

	actualRequest, err := http.NewRequest(request.Method, request.URL.String(), reader)
	if err != nil {
		t.Fatalf("test fixture '%s' not loaded %v", filename, err)
	}
	return actualRequest
}
