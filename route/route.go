package route

import (
	"github.com/zenazn/goji"
	"github.com/tb/simpleIssueTracker/system"
	"github.com/tb/simpleIssueTracker/controller"
)

func PrepareRoutes(application *system.Application) {
	goji.Post("/services/user/signup", application.Route(&controller.Controller{}, true,"SignUp"))
	goji.Post("/services/user/login", application.Route(&controller.Controller{}, true,"Login"))
	goji.Post("/services/user/assign/issue", application.Route(&controller.Controller{}, false,"CreateIssue"))
	goji.Post("/services/user/list/issue", application.Route(&controller.Controller{}, false,"ListIssue"))
	goji.Post("/services/user/delete/issue", application.Route(&controller.Controller{}, false,"DeleteIssue"))
	goji.Post("/services/user/update/issue", application.Route(&controller.Controller{}, false,"UpdateIssue"))

}
