package weoscontroller_test

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	weoscontroller "github.com/wepala/weos-controller"
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
	request, err := weoscontroller.LoadHttpRequestFixture("./_endpoint.input.http")
	if err != nil {
		t.Fatalf("unexpected error loading fixture '%s' '%s'", "./_endpoint.input.http", err)
	}
	var bodyStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(request.Body).Decode(&bodyStruct)
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

func TestAPI_ResponseRecording(t *testing.T) {
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
	}, api.ResponseRecording)
	api.SetEchoInstance(e)
	e.ServeHTTP(rec, req)
	//confirm file is created
	if _, err := os.Stat("./_endpoint.golden.http"); os.IsNotExist(err) {
		t.Fatalf("expected fixture file to be created")
	}

	//confirm contents of file
	response, err := weoscontroller.LoadHttpResponseFixture("./_endpoint.golden.http", req)
	if err != nil {
		t.Fatalf("unexpected error loading fixture '%s' '%s'", "./_endpoint.golden.http", err)
	}
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("error unmarshalling response '%s'", err)
	}

	if string(body) != "test" {
		t.Errorf("expected the name to be '%s', got '%s'", "test", body)
	}

	//delete file
	err = os.Remove("./_endpoint.golden.http")
	if err != nil {
		t.Fatalf("unable to delete test fixture created '%s', got error '%s'", "./_endpoint.golden.http", err)
	}
}
