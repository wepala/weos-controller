package weoscontroller_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	module "github.com/wepala/sellwithwe-module"
	sellwithwemodule "github.com/wepala/sellwithwe-module"
	"github.com/wepala/weos"
	weoscontroller "github.com/wepala/weos-controller"
)

func TestMiddleware_CustomErrorHanlderDomain(t *testing.T) {
	e := echo.New()
	payload := &sellwithwemodule.StorePayload{
		Name: "",
		Phones: []*module.PhoneNumber{
			{
				Value: "123456",
			},
		},
		Emails: []*module.Email{
			{
				Value: "somemail@mail.com",
			},
		},
	}
	reqBytes, _ := json.Marshal(payload)
	body := bytes.NewReader(reqBytes)
	request := httptest.NewRequest("POST", "http://localhost/stores", body)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rw := httptest.NewRecorder()
	c := e.NewContext(request, rw)
	err := &weoscontroller.WeOSControllerError{}
	err.Message = "Some Error"
	err.StatusCode = 0
	err.Err = &weos.DomainError{
		WeOSError: &weos.WeOSError{
			Application: "WEOS",
			AccountID:   "1234",
		},
		EntityID:   "",
		EntityType: "Store",
	}
	weoscontroller.CustomErrorHandler(err, c)
	if c.Response().Status != 400 {
		t.Errorf("expected status code to be 400 with error type DomainError got %d", c.Response().Status)
	}
}

func TestMiddleware_CustomErrorHanlderWeOS(t *testing.T) {
	e := echo.New()
	payload := &sellwithwemodule.StorePayload{
		Name: "",
		Phones: []*module.PhoneNumber{
			{
				Value: "123456",
			},
		},
		Emails: []*module.Email{
			{
				Value: "somemail@mail.com",
			},
		},
	}
	reqBytes, _ := json.Marshal(payload)
	body := bytes.NewReader(reqBytes)
	request := httptest.NewRequest("POST", "http://localhost/stores", body)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rw := httptest.NewRecorder()
	c := e.NewContext(request, rw)
	err := &weoscontroller.WeOSControllerError{}
	err.Message = "Some Error"
	err.StatusCode = 0
	err.Err = &weos.WeOSError{
		Application: "WEOS",
		AccountID:   "1234",
	}
	weoscontroller.CustomErrorHandler(err, c)
	if c.Response().Status != 500 {
		t.Errorf("expected status code to be 500 with error type WeOSError got %d", c.Response().Status)
	}
}
