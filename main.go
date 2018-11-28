package main

import (
	"flag"
	"github.com/tb/simpleIssueTracker/route"
	"github.com/tb/simpleIssueTracker/system"
	"github.com/zenazn/goji"
)

func main()  {
	var application = &system.Application{}
	goji.Use(application.ApplyAuth)
	route.PrepareRoutes(application)
	flag.Set("bind","localhost:8086")
	goji.Serve()

}