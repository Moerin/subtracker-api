package main

import (
	"fmt"

	// Project related packages
	"github.com/yonmey/subtracker-api/lib/dbc"
	"github.com/yonmey/subtracker-api/subscription"

	// Echo framework
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

const serverPort = ":8080"

func init() {
	dbc.DbInit()
}

func main() {
	e := echo.New()

	// Routing
	// Subscription consulting
	e.GET("/subscriptions", subscription.GetAll)
	e.GET("/subscription/:id", subscription.GetOne)

	// CRUD
	e.POST("/subscription/add", subscription.Add)
	e.DELETE("/subscription/delete/:id", subscription.Delete)
	e.PUT("/subscription/update/:id", subscription.Update)

	// Server
	fmt.Println("Serving on port", serverPort)
	e.Run(standard.New(serverPort))
}
