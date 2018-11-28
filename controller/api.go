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

func (controller *Controller) Login(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response, err := services.ValidateUserNameandPassword(data)

	if err != nil {
		return nil, err
	}
	return response, nil
}
func (controller *Controller) CreateIssue(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data model.Issue
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response, err := services.CreateIssue(data,controller.GetUserContext(c))

	if err != nil {
		return nil, err
	}
	return response, nil
}
func (controller *Controller) UpdateIssue(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data model.Issue
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response, err := services.UpdateIssue(data,controller.GetUserContext(c))

	if err != nil {
		return nil, err
	}
	return response, nil
}
func (controller *Controller) DeleteIssue(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	decoder := json.NewDecoder(r.Body)
	var data map[string]int
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	issueId:=data["issueId"]
	response, err := services.DeleteIssue(issueId,controller.GetUserContext(c))

	if err != nil {
		return nil, err
	}
	return response, nil
}
func (controller *Controller) ListIssue(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	response, err := services.ListIssue(controller.GetUserContext(c))

	if err != nil {
		return nil, err
	}
	return response, nil
}