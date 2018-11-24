package route

import (
	"github.com/zenazn/goji"
	"github.com/tb/simpleIssueTracker/system"
	"github.com/tb/simpleIssueTracker/controller"
)

func PrepareRoutes(application *system.Application) {
	goji.Post("/services/user/signup", application.Route(&controller.Controller{}, "SignUp"))
	goji.Post("/services/list/searched/keys", application.Route(&controller.Controller{}, "ListUserSearchInputs"))
	goji.Post("/services/list/image/urls/by/searchedkey", application.Route(&controller.Controller{}, "GetSearchedImageUrlsFromDB"))

}
