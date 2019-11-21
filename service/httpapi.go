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
	"sort"
	"strconv"
	"strings"
)

type mockHandler struct {
	statusCode  int
	contentType string
	content     string
	pathConfig  PathConfig
}

func (h *mockHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//return a response based on the status code set on the handler with the content type header set to the content type
	rw.Header().Add("Content-Type", h.contentType)
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(h.statusCode)
	tmpl, err := template.New("mock").Parse(h.content)
	if err != nil {
		log.Errorf("error rendering mock : '%s'", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(rw, h.pathConfig.Data); err != nil {
		log.Errorf("error rendering mock : '%v'", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func NewMockHandler(statusCode int, content *openapi3.Content) (*mockHandler, error) {
	//check the content type and set the appropriate variable on the handler
	keys := reflect.ValueOf(*content).MapKeys()

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
					log.Debugf("type: %s", exampleString)
					log.Debugf("content-type: %s", contentType)
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
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticFolder))))
	config := service.GetConfig()
	if config != nil {
		paths := make([]string, 0, len(config.Paths))
		for k := range config.Paths {
			paths = append(paths, k)
		}
		sort.Strings(paths)
		for _, path := range paths {
			pathObject := config.Paths[path]
			var pathMethods []string

			for method, operation := range pathObject.Operations() {
				var responseContent *openapi3.Content
				var statusCode int
				var err error
				pathMethods = append(pathMethods, method)
				for statusCodeString, responseRef := range operation.Responses {
					log.Debug(path + " " + statusCodeString + " has mock responses")
					statusCode, err = strconv.Atoi(statusCodeString)
					if err != nil {
						log.Debugf("could not mock the response for the path '%s' for the operation '%s' because the code statusCode %s could not be converted to an integer", path, method, statusCodeString)
					} else {
						responseContent = &responseRef.Value.Content
					}
				}

				if responseContent != nil {
					handler, err := NewMockHandler(statusCode, responseContent)
					if err != nil {
						log.Errorf("could not mock the response for the path '%s' for the operation '%s' because the mock handler could not be created because '%s'", path, method, err)
					}
					router.Handle(path, handler).Methods(method)
				}
			}

			//Add handler for each path's OPTIONS call
			pathMethods = append(pathMethods, "OPTIONS")
			router.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
				//return a response based on the status code set on the handler with the content type header set to the content type
				rw.Header().Add("Access-Control-Allow-Methods", strings.Join(pathMethods, ", "))
				rw.Header().Add("Access-Control-Allow-Origin", "*")
				rw.WriteHeader(200)
			}).Methods("OPTIONS")

		}
	}

	return router
}

func NewHTTPServer(service ServiceInterface, staticFolder string) http.Handler {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticFolder))))
	config := service.GetConfig()

	if config != nil {
		paths := make([]string, 0, len(config.Paths))
		for k := range config.Paths {
			paths = append(paths, k)
		}
		sort.Strings(paths)
		for _, path := range paths {
			pathObject := config.Paths[path]
			var pathMethods []string
			for method, _ := range pathObject.Operations() {
				pathMethods = append(pathMethods, method)
				n := negroni.Classic()
				pathConfig, err := service.GetPathConfig(path, strings.ToLower(method))
				if err != nil {
					log.Errorf("error encountered getting the path config for the route '%s', got: '%s'", path, err.Error())
				}
				handlers, err := service.GetHandlers(pathConfig)
				if err != nil {
					log.Errorf("error encountered retrieving the handlers for the route '%s', got: '%s'", path, err.Error())
				}
				for _, handler := range handlers {
					n.UseHandler(handler)
				}
				router.Handle(path, n).Methods(method)
				log.Debugf("added %d handler(s) to path %s %s", len(handlers), path, method)

				//Add handler for each path's OPTIONS call
				pathMethods = append(pathMethods, "OPTIONS")
				router.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
					//return a response based on the status code set on the handler with the content type header set to the content type
					rw.Header().Add("Access-Control-Allow-Methods", strings.Join(pathMethods, ", "))
					rw.Header().Add("Access-Control-Allow-Origin", "*")
					rw.Header().Add("Accept", "text/html,application/xhtml+xml,application/json;q=0.9,*/*;q=0.8")
					rw.WriteHeader(200)
				}).Methods("OPTIONS")
			}

		}
	}

	return router
}

func GenerateStaticPages(serviceInterface ServiceInterface, route string, data []interface{}) {

}
