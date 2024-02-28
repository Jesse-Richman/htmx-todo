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

func (h TodoHandler) Main(c echo.Context) error {
	return render(c, view.Main(h.DB.GetTodos()))
}

func (h TodoHandler) Table(c echo.Context) error {
    // filter := c.QueryParam("filter")
    // parts := strings.Split(c.Path(), "/")
    parts := strings.Split(c.Request().URL.Path, "/")
    filter := parts[len(parts)-1]
    log.Infof("filtering with %s", filter)
    var todos []model.Todo

    if filter != "" {
        state := false
        if filter == "complete" {
            state = true
        }
        todos = h.DB.GetTodosByDone(state)
    } else {
        todos = h.DB.GetTodos()
    }

	return render(c, view.Table(todos))
}
func (h TodoHandler) Create(c echo.Context) error {
	desc := c.FormValue("description")
    todo := h.DB.SaveTodo(desc)
    c.Response().Header().Add("HX-Trigger", "editTodo")
	return render(c, view.ViewTodo(todo))
}

func (h TodoHandler) Delete(c echo.Context) error {
	id := c.Param("id")
    h.DB.DeleteTodo(id)
    c.Response().Header().Add("HX-Trigger", "editTodo")
	return nil
}

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
    if doneStr == "on" {
        done = true
    } else {
        done = !todo.Done
    }

    todo = h.DB.UpdateTodo(id, desc, done)

    c.Response().Header().Add("HX-Trigger", "editTodo")
	return render(c, view.ViewTodo(todo))
}

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

func (h TodoHandler) Metrics(c echo.Context) error {
    var count int
    for _, todo := range h.DB.GetTodos() {
        if !todo.Done {
            count++
        }
    }
    return render(c, view.Metrics(count))
}

func (h TodoHandler) Clear(c echo.Context) error {
    todos := h.DB.GetTodos()
    for _,todo := range todos {
        if todo.Done {
            h.DB.DeleteTodo(todo.ID)
        }
    }
    return render(c, view.Table(h.DB.GetTodos()))
}

