package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"jesse.richman/todo/handler"
	"jesse.richman/todo/model"
)

func main() {
	app := echo.New()
    app.Use(middleware.Logger())
    
    // mock database
    db := model.DB{}
    db.InitTestData()

    // setup handlers
	todoHandler := handler.TodoHandler{DB: &db}

	app.GET("/", todoHandler.Main)
	app.POST("/todo", todoHandler.Create)
    app.DELETE("/todo/:id", todoHandler.Delete)
    app.PATCH("/todo/:id", todoHandler.Update)
    app.GET("/todo/:id", todoHandler.GetTodo)

    app.GET("/filter/*", todoHandler.List)
    app.GET("/todo/metrics", todoHandler.Metrics)
    app.DELETE("/todo/clear", todoHandler.Clear)

    // start the app
	app.Logger.Fatal(app.Start(":3000"))
}
