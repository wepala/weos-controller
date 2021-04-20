package weoscontroller_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	weoscontroller "github.com/wepala/weos-controller"
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

func TestAPI_Authenticate(t *testing.T) {

	t.Run("successful jwt authenticate", func(t *testing.T) {
		// Setup
		e := echo.New()
		key := "secureSecretText"
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{})
		signedToken, err := token.SignedString([]byte(key))
		var bearer = "Bearer " + signedToken
		if err != nil {
			t.Errorf("got an error setting up tests %s", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/endpoint", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("Authorization", bearer)
		rec := httptest.NewRecorder()
		api := &weoscontroller.API{Config: &weoscontroller.APIConfig{
			RecordingBaseFolder: ".",
			JWTConfig: &weoscontroller.JWTConfig{
				Key:             key,
				SigningKeys:     map[string]interface{}{},
				Certificate:     nil,
				CertificatePath: "",
				TokenLookup:     "",
				AuthScheme:      "",
				SigningMethod:   "HS256",
				ContextKey:      "",
			},
		}}

		e.POST("/endpoint", func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}, api.Authenticate)

		api.SetEchoInstance(e)
		e.ServeHTTP(rec, req)

		response := rec.Result()
		defer response.Body.Close()

		if response.StatusCode != 200 {
			t.Errorf("expected the status code to be %d, got %d", 200, response.StatusCode)
		}
	})
	t.Run("jwt authenticate no token", func(t *testing.T) {
		// Setup
		e := echo.New()
		key := "secureSecretText"
		req := httptest.NewRequest(http.MethodPost, "/endpoint", strings.NewReader(`{"name":"Sojourner Truth","email":"sojourner@examle.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		api := &weoscontroller.API{Config: &weoscontroller.APIConfig{
			RecordingBaseFolder: ".",
			JWTConfig: &weoscontroller.JWTConfig{
				Key:             key,
				SigningKeys:     map[string]interface{}{},
				Certificate:     nil,
				CertificatePath: "",
				TokenLookup:     "",
				AuthScheme:      "",
				SigningMethod:   "HS256",
				ContextKey:      "",
			},
		}}

		e.POST("/endpoint", func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}, api.Authenticate)

		api.SetEchoInstance(e)
		e.ServeHTTP(rec, req)

		response := rec.Result()
		defer response.Body.Close()

		if response.StatusCode != 400 {
			t.Errorf("expected the status code to be %d, got %d", 400, response.StatusCode)
		}

		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("got an error getting response body %s", err)
		}
		responseMessage := strings.TrimSpace(string(bodyBytes))
		if responseMessage != `{"message":"missing or malformed jwt"}` {
			t.Errorf("expected the response message to be %s got %s", `{"message":"missing or malformed jwt"}`, responseMessage)
		}
	})
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
