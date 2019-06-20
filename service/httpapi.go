package service

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

type mockHandler struct {
	statusCode  int
	contentType string
	content     string
}

func (*mockHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	//TODO return a response based on the status code set on the handler with the content type header set to the content type
	panic("implement me")
}

func NewMockHandler(statusCode int, content openapi3.Content) (*mockHandler, error) {
	//TODO check the content type and set the appropriate variable on the handler
	return &mockHandler{
		statusCode: statusCode,
	}, nil
}

func NewMockHTTPServer(service ServiceInterface, staticFolder string) http.Handler {
	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(negroni.New(negroni.NewStatic(http.Dir(staticFolder))))
	n := negroni.Classic()
	config := service.GetConfig()

	if config != nil {
		for path, pathObject := range config.ApiConfig.Paths {
			for method, operation := range pathObject.Operations() {
				for statusCodeString, responseRef := range operation.Responses {
					statusCode, err := strconv.Atoi(statusCodeString)
					if err != nil {
						log.Debugf("could not mock the response for the path '%s' for the operation '%s' because the code statusCode %s could not be converted to an integer", path, method, statusCodeString)
					} else {
						handler, err := NewMockHandler(statusCode, responseRef.Value.Content)
						if err != nil {
							log.Debugf("could not mock the response for the path '%s' for the operation '%s' because the mock handler could not be created", path, method)
						}
						router.Handle(path, handler)
					}
				}

			}

		}
	}

	n.UseHandler(router)
	return n
}

func NewHTTPServer(service ServiceInterface, staticFolder string) http.Handler {
	//TODO setup a handler using gorilla + negroni (or just negroni?)
	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(negroni.New(negroni.NewStatic(http.Dir(staticFolder))))
	n := negroni.Classic()
	config := service.GetConfig()

	if config != nil {
		for path, pathObject := range config.ApiConfig.Paths {
			for method, operation := range pathObject.Operations() {
				for statusCodeString, responseRef := range operation.Responses {
					statusCode, err := strconv.Atoi(statusCodeString)
					if err != nil {
						log.Debugf("could not mock the response for the path '%s' for the operation '%s' because the code statusCode %s could not be converted to an integer", path, method, statusCodeString)
					} else {
						handler, err := NewMockHandler(statusCode, responseRef.Value.Content)
						if err != nil {
							log.Debugf("could not mock the response for the path '%s' for the operation '%s' because the mock handler could not be created", path, method)
						}
						router.Handle(path, handler)
					}
				}

			}

		}
	}
	//TODO add middleware that returns the html response
	n.UseHandler(router)
	return n
}
