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
	var mockStatusVal int
	ok := false
	showStatusCodeError := false
	showExampleError := false

	mockStatusCode := r.Header.Get("X-Mock-Status-Code")
	mockExample := r.Header.Get("X-Mock-Example")

	var err error

	if mockStatusCode != "" {
		showStatusCodeError = true
		mockStatusVal, err = strconv.Atoi(mockStatusCode)
		if err != nil {
			log.Errorf("Error converting string to integer: %s", err.Error())
		}
	}

	if mockExample != "" {
		showExampleError = true
	}

	for _, operation := range h.PathInfo.Operations() {
		var responseContent *openapi3.Content

		if showStatusCodeError {
			for statusCodeString, responseRef := range operation.Responses {
				responseContent = &responseRef.Value.Content
				if statusCodeString == mockStatusCode {
					ok = h.getMockResponses(responseContent, rw, r)
				}
			}
		} else if !showStatusCodeError {
			var defaultResponse *openapi3.ResponseRef
			if operation.Responses.Default() != nil {
				defaultResponse = operation.Responses.Default()
			} else {
				if operation.Responses.Get(200) != nil {
					defaultResponse = operation.Responses.Get(200)
				} else {
					rw.Write([]byte("Error: Cannot mock this endpoint"))
					return
				}
			}
			ok = h.getMockResponses(&defaultResponse.Value.Content, rw, r)
		}
		if ok {
			break
		}
	}
		//rw.Header().Add("Access-Control-Allow-Origin", "*")
		//rw.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
	if !ok {
		rw.Header().Add("Content-Type", "text/plain")
		if showExampleError {
			rw.WriteHeader(mockStatusVal)
			rw.Write([]byte("There is no mocked response with example named " + mockExample))
			return
		} else if showStatusCodeError {
			rw.WriteHeader(200)
			rw.Write([]byte("There is no mocked response for status code " + mockStatusCode))
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

func(h *MockHandler) getMockResponses (responseContent *openapi3.Content, rw http.ResponseWriter, r *http.Request) bool{
	var mockStatusVal int
	var mockExampleLengthVal int

	showStatusCodeError := false
	showExampleError := false
	showContentType := false

	mockStatusCode := r.Header.Get("X-Mock-Status-Code")
	mockContentType := r.Header.Get("X-Mock-Content-Type")
	mockExample := r.Header.Get("X-Mock-Example")
	mockExampleLength := r.Header.Get("X-Mock-Example-Length")

	var err error

	if mockStatusCode != "" {
		showStatusCodeError = true
		mockStatusVal, err = strconv.Atoi(mockStatusCode)
		if err != nil {
			log.Errorf("Error converting string to integer: %s", err.Error())
		}
	}

	if mockExampleLength != "" {
		mockExampleLengthVal, err = strconv.Atoi(mockExampleLength)
		if err != nil {
			log.Errorf("Error converting string to integer: %s", err.Error())
		}
	}

	if mockContentType != "" {
		showContentType = true
	}

	if mockExample != "" {
		showExampleError = true
	}

	keys := reflect.ValueOf(*responseContent).MapKeys()
	if len(keys) > 0 {
		for _, key := range keys {
			contentType := key.String()
			if (showContentType && contentType == mockContentType) || (!showContentType) {
				var c *openapi3.MediaType
				if h.PathInfo.GetOperation("OPTIONS") != nil {
					if h.PathInfo.GetOperation("OPTIONS").Responses.Get(mockStatusVal) != nil {
						rw.Header().Add("Access-Control-Allow-Origin", h.PathInfo.GetOperation("OPTIONS").Responses.Get(mockStatusVal).Value.Headers["Access-Control-Allow-Origin"].Value.Schema.Value.Example.(string))
						rw.Header().Add("Access-Control-Allow-Headers", h.PathInfo.GetOperation("OPTIONS").Responses.Get(mockStatusVal).Value.Headers["Access-Control-Allow-Headers"].Value.Schema.Value.Example.(string))
					} else {
						rw.Header().Add("Access-Control-Allow-Origin", "*")
						rw.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
					}
				} else {
					rw.Header().Add("Access-Control-Allow-Origin", "")
					rw.Header().Add("Access-Control-Allow-Headers", "")
				}

				if showContentType {
					rw.Header().Add("Content-Type", mockContentType)
					c = responseContent.Get(mockContentType)
				} else {
					rw.Header().Add("Content-Type", contentType)
					c = responseContent.Get(contentType)
				}

				if showStatusCodeError {
					rw.WriteHeader(mockStatusVal)
				} else {
					rw.WriteHeader(200)
				}

				if c != nil && (c.Example != nil || c.Examples != nil) {
					if c.Example != nil {
						switch x := c.Example.(type) {
						case string:
							log.Infof("type: %s", x)
							rw.Write([]byte(c.Example.(string)))
							return true
						default:
							if c.Extensions["example"] != nil {
								//found that the Extensions property was a better way to access the raw data
								example := c.Extensions["example"].(json.RawMessage)
								exampleString, err := example.MarshalJSON()
								if err != nil {
									log.Errorf("Error marshalling json: %s", err.Error())
									return false
								}

								//example := string(data)[11:len(string(data))1]
								log.Infof("type: %s", exampleString)
								log.Infof("contenttype: %s", contentType)
								rw.Write(exampleString)
								return true
							}
						}
					} else if c.Examples != nil {
						if showExampleError {
							for name, example := range c.Examples {
								if name == mockExample {
									rw.Write([]byte(example.Value.Value.(string)))
									return true
								}
							}
						} else {
							rw.Write([]byte("There are multiple examples defined. Please specify one using the X-Mock-Example header"))
							return true
						}
					}
				}
				if contentType == "application/json" {
					if c.Schema.Value.Example != nil {
						body, err := json.Marshal(c.Schema.Value.Example)
						if err != nil {
							log.Errorf("Error mashalling json, %q", err.Error())
							return false
						}
						rw.Write(body)
					} else if c.Schema.Value.Items.Value.Example != nil {
						arrayLength := mockExampleLengthVal

						exampleValue := c.Schema.Value.Items.Value.Example
						exampleArray := make([]interface{}, arrayLength)
						exampleArray[0] = exampleValue
						body, err := json.Marshal(exampleArray)
						//fmt.Println(c.Schema.Value.Type)
						if err != nil {
							log.Errorf("Error mashalling json, %q", err.Error())
							return false
						}
						rw.Write(body)
					}
					return true
				}
			}
		}
	}
	return false
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
				handlers, err := service.GetHandlers(pathConfig, &MockHandler{PathInfo: pathObject})
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
