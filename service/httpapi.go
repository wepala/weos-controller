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

	//this is a flag to check if the returned response is written or not
	ok := false

	//values pulled from headers are stored here
	var mockStatusVal int

	//flags for headers are set here
	var showStatusCodeError = false
	var showExampleError = false

	//values from headers are pulled
	mockStatusCode := r.Header.Get("X-Mock-Status-Code")
	mockExample := r.Header.Get("X-Mock-Example")

	var err error

	if mockStatusCode != "" {
		showStatusCodeError = true

		//convert status code from string to integer
		mockStatusVal, err = strconv.Atoi(mockStatusCode)
		if err != nil {
			log.Errorf("Error converting string to integer: %s", err.Error())
		}
	}

	if mockExample != "" {
		showExampleError = true
	}

	//for each operation that has been retrieved from the swagger file
	for _, operation := range h.PathInfo.Operations() {
		var responseReference *openapi3.ResponseRef

		//if a status code was retrieved from the header
		if showStatusCodeError {
			// for the corresponding status codes and responses from the swagger file
			for statusCodeString, responseRef := range operation.Responses {
				responseReference = responseRef
				//if the two status codes match
				if statusCodeString == mockStatusCode {
					ok = h.getMockResponses(responseReference, rw, r)
				}
				if ok{
					break
				}
			}
			//else if a status code was not retrieved from the header
		} else if !showStatusCodeError {
			var defaultResponse *openapi3.ResponseRef
			// Check if the responses contain a default response
			if operation.Responses.Default() != nil {
				defaultResponse = operation.Responses.Default()
			} else {
				//check if there's a 200 response
				if operation.Responses.Get(200) != nil {
					defaultResponse = operation.Responses.Get(200)
				} else {
					rw.Write([]byte("Error: Cannot mock this endpoint"))
					return
				}
			}
			ok = h.getMockResponses(defaultResponse, rw, r)
		}
		if ok {
			break
		}
	}
	//if no example was written throught this process, we write one of these outputs
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
}

func(h *MockHandler) getMockResponses (responseRef *openapi3.ResponseRef, rw http.ResponseWriter, r *http.Request) bool{

	//values pulled from headers are stored here
	var mockStatusVal int
	var mockExampleLengthVal int

	//flags for headers are set here
	showStatusCodeError := false
	showExampleError := false
	showContentType := false

	//values from headers are pulled
	mockStatusCode := r.Header.Get("X-Mock-Status-Code")
	mockContentType := r.Header.Get("X-Mock-Content-Type")
	mockExample := r.Header.Get("X-Mock-Example")
	mockExampleLength := r.Header.Get("X-Mock-Example-Length")

	var err error

	if mockStatusCode != "" {
		showStatusCodeError = true
		//convert status code from string to integer
		mockStatusVal, err = strconv.Atoi(mockStatusCode)
		if err != nil {
			log.Errorf("Error converting string to integer: %s", err.Error())
		}
	}

	if mockExampleLength != "" {
		//convert example length from string to integer
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

	responseContent := &responseRef.Value.Content

	//store all the kets from the content
	keys := reflect.ValueOf(*responseContent).MapKeys()


	//attach headers
	if responseRef.Value.Headers != nil {
		for key, headerVal := range responseRef.Value.Headers{
			rw.Header().Add(key, headerVal.Value.Schema.Value.Example.(string))
		}
	}else{
		rw.Header().Add("Access-Control-Allow-Origin", "")
		rw.Header().Add("Access-Control-Allow-Headers", "")
	}

	//if there is at least 1 key, we start the process
	if (len(keys) == 1) || (len(keys) > 1 && showContentType) {
		for _, key := range keys {
			//retrieve the content type
			contentType := key.String()
			//if a content type was pulled from the headers and it matches the current content type from the keys
			if (showContentType && contentType == mockContentType) || (!showContentType) {
				var c *openapi3.MediaType

				//if a content type was pulled from the headers, set it here, otherwise use the one from the key
					if showContentType {
					rw.Header().Add("Content-Type", mockContentType)
					c = responseContent.Get(mockContentType)
				} else {
					rw.Header().Add("Content-Type", contentType)
					c = responseContent.Get(contentType)
				}

				//if a status was pulled from the headers, set it here, otherwise use 200
				if showStatusCodeError {
					rw.WriteHeader(mockStatusVal)
				} else {
					rw.WriteHeader(200)
				}

				//so long as the content is not empty and there is at least 1 example
				if c != nil && (c.Example != nil || c.Examples != nil) {
					if c.Example != nil {
						switch x := c.Example.(type) {
						case string:
							//write the example if it is a simple string
							log.Infof("type: %s", x)
							rw.Write([]byte(c.Example.(string)))
							return true
						default:
							//if the example is not a string, it will do the following
							if c.Extensions["example"] != nil {
								//found that the Extensions property was a better way to access the raw data
								example := c.Extensions["example"].(json.RawMessage)
								exampleString, err := example.MarshalJSON()
								if err != nil {
									log.Errorf("Error marshalling json: %s", err.Error())
									return false
								}

								log.Infof("type: %s", exampleString)
								log.Infof("contenttype: %s", contentType)
								rw.Write(exampleString)
								return true
							}
						}
						//else if there are multiple examples
					} else if c.Examples != nil {
						//if the example name was retrieved from the header
						if showExampleError {
							for name, example := range c.Examples {
								if name == mockExample {
									//write that specified example
									rw.Write([]byte(example.Value.Value.(string)))
									return true
								}
							}
						} else {
							//else write error
							rw.Write([]byte("There are multiple examples defined. Please specify one using the X-Mock-Example header"))
							return true
						}
					}
				}
				//if the content is json it'll go through this process
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
						if err != nil {
							log.Errorf("Error mashalling json, %q", err.Error())
							return false
						}
						rw.Write(body)
					}
					return true
				}
				//if there is no content type that was pulled from the headers
			}
		}
	}else if len(keys) > 1 && !showContentType{
		rw.Write([]byte("There are multiple content types defined. Please specify one using the X-Mock-Content-Type header"))
		return true
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
