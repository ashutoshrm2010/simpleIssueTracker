package controller

import (
	"github.com/tb/simpleIssueTracker/system"
	"net/http"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"github.com/tb/simpleIssueTracker/model"
	"fmt"
	"github.com/tb/simpleIssueTracker/services"
)

type Controller struct {
	system.Controller
}

func (controller *Controller) SignUp(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data model.User
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("error ",err)
		return nil, err
	}
	fmt.Println("data ",data)
	response, err := services.SignUp(data)

	if err != nil {
		return nil, err
	}
	return response, nil
}