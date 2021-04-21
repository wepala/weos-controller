package weoscontroller_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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
	t.Run("successful jwt authenticate certificate", func(t *testing.T) {
		// Setup
		e := echo.New()

		privatekey, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			t.Errorf("got an error setting up tests %s", err)
		}
		publicKey := &privatekey.PublicKey
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			t.Errorf("got an error setting up tests %s", err)
		}
		publicKeyBlock := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
		publicPem, err := os.Create("./fixtures/certificates/publicRS.txt")
		if err != nil {
			t.Errorf("got an error setting up tests %s", err)
		}
		err = pem.Encode(publicPem, publicKeyBlock)
		if err != nil {
			t.Errorf("got an error setting up tests %s", err)
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{})
		signedToken, err := token.SignedString(privatekey)
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
				Key:             "",
				SigningKeys:     map[string]interface{}{},
				Certificate:     nil,
				CertificatePath: "./fixtures/certificates/publicRS.txt",
				TokenLookup:     "",
				AuthScheme:      "",
				SigningMethod:   "RS256",
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
	t.Run("jwt authenticate expired token", func(t *testing.T) {
		// Setup
		e := echo.New()
		key := "secureSecretText"
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: jwt.TimeFunc().Unix() - 100})
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

		if response.StatusCode != 401 {
			t.Errorf("expected the status code to be %d, got %d", 401, response.StatusCode)
		}
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("got an error getting response body %s", err)
		}
		responseMessage := strings.TrimSpace(string(bodyBytes))
		if responseMessage != `{"message":"invalid or expired jwt"}` {
			t.Errorf("expected the response message to be %s got %s", `{"message":"invalid or expired jwt"}`, responseMessage)
		}
	})
	t.Run("jwt authenticate invalid token", func(t *testing.T) {
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
				Key:             "different",
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

		if response.StatusCode != 401 {
			t.Errorf("expected the status code to be %d, got %d", 401, response.StatusCode)
		}
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("got an error getting response body %s", err)
		}
		responseMessage := strings.TrimSpace(string(bodyBytes))
		if responseMessage != `{"message":"invalid or expired jwt"}` {
			t.Errorf("expected the response message to be %s got %s", `{"message":"invalid or expired jwt"}`, responseMessage)
		}
	})
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
