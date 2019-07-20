package service

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
)

type mockHandler struct {
	statusCode  int
	contentType string
	content     string
	pathConfig  PathConfig
}

func (h *mockHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//return a response based on the status code set on the handler with the content type header set to the content type
	rw.WriteHeader(h.statusCode)
	rw.Header().Set("Content-Type", h.contentType)
	tmpl, err := template.New("mock").Parse(h.content)
	if err != nil {
		log.Errorf("error rendering mock : '%v'", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(rw, h.pathConfig.Data); err != nil {
		log.Errorf("error rendering mock : '%v'", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func NewMockHandler(statusCode int, content openapi3.Content) (*mockHandler, error) {
	//check the content type and set the appropriate variable on the handler
	keys := reflect.ValueOf(content).MapKeys()

	if len(keys) > 0 {
		contentType := keys[0].String()
		c := content.Get(contentType)
		if c != nil && c.Example != nil {

			switch x := c.Example.(type) {
			case string:
				log.Infof("type: %s", x)
				return &mockHandler{
					statusCode: statusCode,
					content:    c.Example.(string),
				}, nil
			default:
				if c.Extensions["example"] != nil {
					//found that the Extensions property was a better way to access the raw data
					example := c.Extensions["example"].(json.RawMessage)
					exampleString, err := example.MarshalJSON()
					if err != nil {
						return nil, err
					}
					//trim {"example": from the front and "}" from the end

					//example := string(data)[11:len(string(data))-1]
					log.Infof("type: %s", exampleString)
					return &mockHandler{
						statusCode:  statusCode,
						content:     string(exampleString),
						contentType: contentType,
					}, nil
				}
			}

		}
	}

	return &mockHandler{
		statusCode: statusCode,
		content:    "This endpoint was not mocked",
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
							log.Errorf("could not mock the response for the path '%s' for the operation '%s' because the mock handler could not be created because '%s'", path, method, err)
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
