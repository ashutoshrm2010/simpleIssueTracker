package system

import (
	"net/http"
	"encoding/json"
	"reflect"
	"github.com/zenazn/goji/web"
	"errors"
)
type UserContext struct {
	UserId   int
	Name string
	UserName string
	EmailId  string
}
type Application struct {
}
type Controller struct {
}

func (application *Application) Route(controller interface{},isPublic bool, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		if c.Env["AuthFailed"].(bool)&&!isPublic {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			response := make(map[string]interface{})
			response["message"] = errors.New("Unauthorized request")
			errResponse, _ := json.Marshal(response)
			w.Write(errResponse)
		} else {
				methodValue := reflect.ValueOf(controller).MethodByName(route)
				methodInterface := methodValue.Interface()

				method := methodInterface.(func(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error))
				result, err := method(c, w, r)

				if c.Env["Content-Type"] != nil {
					w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
				} else {
					w.Header().Set("Content-Type", "application/json")
				}

				if (err != nil) {
					response := make(map[string]interface{})
					response["message"] = err.Error()
					errResponse, _ := json.Marshal(response)
					w.Write(errResponse)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(result)
				}

			}

	}
	return fn
}

func (controller *Controller) GetUserContext(c web.C) *UserContext {
	userContext := c.Env["userContext"].(*UserContext)

	return userContext

}
