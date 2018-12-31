package servers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// JSONResponse -
type JSONResponse struct {
	HasError     bool        `json:"hasError"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// JSONRequestHandler -
//type JSONRequestHandler func(data interface{}) JsonResponse
type JSONRequestHandler func(data []byte) JSONResponse

// JSONHandler -
type JSONHandler struct {
	Route      string
	Handler    JSONRequestHandler
	JSONObject interface{}
}

// JSONServer -
type JSONServer struct {
	handlers       []JSONHandler
	enableCors     bool
	routeToHandler map[string]JSONHandler
	lowLevelServer *LowLevelServer
}

// NewJSONServer -
func NewJSONServer() *JSONServer {
	return &JSONServer{
		handlers:       []JSONHandler{},
		enableCors:     false,
		routeToHandler: map[string]JSONHandler{},
		lowLevelServer: NewLowLevelServer(),
	}
}

// SetHandlers -
func (thisRef *JSONServer) SetHandlers(handlers []JSONHandler) {
	thisRef.handlers = handlers
}

// EnableCORS -
func (thisRef *JSONServer) EnableCORS() {
	thisRef.lowLevelServer.EnableCORS()
}

// Run - Server interface
func (thisRef *JSONServer) Run(ipPort string) {

	var lowLevelRequestHelper = func(rw http.ResponseWriter, r *http.Request) {
		var jsonHandler = thisRef.routeToHandler[r.URL.Path]

		// Pass Object
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(rw, "can't read body", http.StatusBadRequest)
			return
		}

		var jsonResponse = jsonHandler.Handler(body)
		json.NewEncoder(rw).Encode(jsonResponse)
	}

	var lowLevelHandlers = []LowLevelHandler{}

	for _, handler := range thisRef.handlers {
		thisRef.routeToHandler[handler.Route] = handler

		lowLevelHandlers = append(lowLevelHandlers, LowLevelHandler{
			Route:   handler.Route,
			Handler: lowLevelRequestHelper,
			Verb:    "POST",
		})
	}

	thisRef.lowLevelServer.SetHandlers(lowLevelHandlers)
	thisRef.lowLevelServer.Run(ipPort)
}

// func (jsonData *JsonData) ToObject(objectInstance interface{}) {
// 	// Do JSON to Object Mapping
// 	objectValue := reflect.ValueOf(objectInstance).Elem()
// 	for i := 0; i < objectValue.NumField(); i++ {
// 		field := objectValue.Field(i)
// 		fieldName := objectValue.Type().Field(i).Name

// 		if valueToCopy, ok := (*jsonData)[fieldName]; ok {
// 			if !field.CanInterface() {
// 				continue
// 			}
// 			switch field.Interface().(type) {
// 			case string:
// 				valueToCopyAsString := reflect.ValueOf(valueToCopy).String()
// 				objectValue.Field(i).SetString(valueToCopyAsString)
// 				break
// 			case int:
// 				valueToCopyAsInt := int64(reflect.ValueOf(valueToCopy).Float())
// 				objectValue.Field(i).SetInt(valueToCopyAsInt)
// 				break
// 			case float64:
// 				valueToCopyAsFloat := reflect.ValueOf(valueToCopy).Float()
// 				objectValue.Field(i).SetFloat(valueToCopyAsFloat)
// 				break
// 			default:
// 			}
// 		}
// 	}
// }

// Get JSON fields
//var jsonData JsonData
//_ = json.NewDecoder(r.Body).Decode(&jsonData)

// TRACE
// if false {
// 	reqAsJSON, _ := json.Marshal(req)
// 	fmt.Println(fmt.Sprintf("%s -> %s", Utils.CallStack(), string(reqAsJSON)))
// }

//jsonData.ToObject(jsonHandler.JsonObject)

// Pass Object
//var response JsonResponse = jsonHandler.Handler(jsonData)