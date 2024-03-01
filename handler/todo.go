package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"jesse.richman/todo/model"
	"jesse.richman/todo/view"
)

type TodoHandler struct {
	DB *model.DB
}

// Returns the main todo view.
func (h TodoHandler) Main(c echo.Context) error {
	return render(c, view.Main(h.DB.GetTodos()))
}

// Return list view of todos and possibly filter based on request path.
func (h TodoHandler) List(c echo.Context) error {
    parts := strings.Split(c.Request().URL.Path, "/")
    filter := parts[len(parts)-1]
    log.Infof("filtering with %s", filter)
    var todos []model.Todo

    if filter != "all" {
        state := false
        if filter == "complete" {
            state = true
        }
        todos = h.DB.GetTodosByDone(state)
    } else {
        todos = h.DB.GetTodos()
    }

	return render(c, view.List(todos))
}

// Create a new todo and add it to the DB. Return todo view.
func (h TodoHandler) Create(c echo.Context) error {
	desc := c.FormValue("description")
    todo := h.DB.SaveTodo(desc)
    c.Response().Header().Add("HX-Trigger", "editTodo")
	return render(c, view.ViewTodo(todo))
}

// Delete todo from the DB and return nil
func (h TodoHandler) Delete(c echo.Context) error {
	id := c.Param("id")
    h.DB.DeleteTodo(id)
    c.Response().Header().Add("HX-Trigger", "editTodo")
	return nil
}

// Update todo's description and/or done state. Return todo view.
func (h TodoHandler) Update(c echo.Context) error {
	id := c.Param("id")
    doneStr := c.FormValue("done")
    desc := c.FormValue("description")
    todo := h.DB.GetTodo(id)
    log.Infof("Updating %s with values: '%s', %s", id, desc, doneStr)

    // do some validation and converting
    if desc == "" {
        desc = todo.Description
    }

    var done bool
    if doneStr == "on" || doneStr == "true" {
        done = true
    } else {
        done = false
    }

    todo = h.DB.UpdateTodo(id, desc, done)

    c.Response().Header().Add("HX-Trigger", "editTodo")
	return render(c, view.ViewTodo(todo))
}

// Get todo by param id and render either the normal or edit view based on
// the 'editMode' query param.
func (h TodoHandler) GetTodo(c echo.Context) error {
	id := c.Param("id")
    editMode := c.QueryParams().Has("editMode")
    todo := h.DB.GetTodo(id)

    if editMode {
        return render(c, view.EditTodo(todo))
    } else {
        return render(c, view.ViewTodo(todo))
    }
}

// Get the remaining and total todos and return metrics view
func (h TodoHandler) Metrics(c echo.Context) error {
    remaining := len(h.DB.GetTodosByDone(false))
    total := len(h.DB.GetTodos())
    return render(c, view.Metrics(remaining, total))
}

// Clear all completed todos from DB and return list view
func (h TodoHandler) Clear(c echo.Context) error {
    for _,todo := range h.DB.GetTodos() {
        if todo.Done {
            h.DB.DeleteTodo(todo.ID)
        }
    }
    return render(c, view.List(h.DB.GetTodos()))
}

