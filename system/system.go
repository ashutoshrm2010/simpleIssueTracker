package system

import (
	"net/http"
	"encoding/json"
	"reflect"
	"github.com/zenazn/goji/web"
)

type Application struct {
}
type Controller struct {
}

func (application *Application) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C,w http.ResponseWriter, r *http.Request) {
		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error))
		result, err := method(c,w, r)
		w.Header().Set("Content-Type", "application/json")

		if (err != nil) {
			response := make(map[string]interface{})
			{
				response["message"] = "something went wrong"
				w.WriteHeader(http.StatusInternalServerError)
			}

			errResponse, _ := json.Marshal(response)
			w.Write(errResponse)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}

	}
	return fn
}