package service

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type MockHandler struct {
	PathInfo *openapi3.PathItem
}

func (h *MockHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//return a response based on the status code set on the handler with the content type header set to the content type
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Content-Type", "text/html")
	mockStausCode := r.Header.Get("X-Mock-Status-Code")
	if mockStausCode != ""{
		mockStatusCode, err := strconv.Atoi(mockStausCode)
		if err != nil{
			log.Errorf("Error converting string to integer: %s", err.Error())
		}
		rw.WriteHeader(mockStatusCode)
	}

	for method, operation := range h.PathInfo.Operations() {
		var responseContent *openapi3.Content
		var err error

		for statusCodeString, responseRef := range operation.Responses {
			_, err = strconv.Atoi(statusCodeString)
			if err != nil {
				log.Debugf("could not mock the response for the path for the operation '%s' because the code statusCode %s could not be converted to an integer", method, statusCodeString)
			}
			responseContent = &responseRef.Value.Content
			if statusCodeString != mockStausCode{
				continue
			}

			keys := reflect.ValueOf(*responseContent).MapKeys()

			if len(keys) > 0 {
				contentType := keys[0].String()
				c := responseContent.Get(contentType)
				if c != nil && (c.Example != nil || c.Examples != nil) {
					if c.Example != nil {
						switch x := c.Example.(type) {
						case string:
							log.Infof("type: %s", x)
							rw.Write([]byte(c.Example.(string)))
							return
						default:
							if c.Extensions["example"] != nil {
								//found that the Extensions property was a better way to access the raw data
								example := c.Extensions["example"].(json.RawMessage)
								exampleString, err := example.MarshalJSON()
								if err != nil {
									log.Errorf("Error marshalling json: %s", err.Error())
									return
								}
								//trim {"example": from the front and "}" from the end

								//example := string(data)[11:len(string(data))1]
								log.Debugf("type: %s", exampleString)
								log.Debugf("contenttype: %s", contentType)
								rw.Write(exampleString)
								return
							}
						}
					}
				}
			}
			rw.Write([]byte("This endpoint is not mocked"))
			return
		}
	}

	//tmpl, err := template.New("mock").Parse(h.content)
	//if err != nil {
	//	log.Errorf("error rendering mock : '%s'", err)
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//}
	//if err := tmpl.Execute(rw, h.pathConfig.Data); err != nil {
	//	log.Errorf("error rendering mock : '%v'", err)
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//}
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
				handlers, err := service.GetHandlers(path, pathConfig, pathObject)
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
					rw.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					rw.WriteHeader(200)
				}).Methods("OPTIONS")
			}

		}
	}

	return router
}

func GenerateStaticPages(serviceInterface ServiceInterface, route string, data []interface{}) {

}
