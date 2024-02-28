package model

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type DB struct {
	todos []Todo
}

func (db *DB) InitTestData() {
    db.todos = []Todo {
        {ID: uuid.NewString(), Description: "program more", Done: false},
        {ID: uuid.NewString(), Description: "finish this", Done: false},
    }
}

func (db *DB) SaveTodo(description string) Todo {
	todo := Todo{ID: uuid.NewString(), Description: description, Done: false}
	db.todos = append([]Todo{todo}, db.todos...)
    
    log.Infof("creating new todo: %v", todo)
	return todo
}

func (db *DB) DeleteTodo(id string) {
	for i, todo := range db.todos {
		if todo.ID == id {
			newTodos := make([]Todo, 0)
			newTodos = append(newTodos, db.todos[:i]...)
			newTodos = append(newTodos, db.todos[i+1:]...)
			db.todos = newTodos
		}
	}
}

func (db *DB) UpdateTodo(id string, desc string, done bool) Todo {
	for i, todo := range db.todos {
		if todo.ID == id {
			db.todos[i].Done = done
			db.todos[i].Description = desc
			return db.todos[i]
		}
	}

	return Todo{Description: desc, Done: done}
}

func (db DB) GetTodo(id string) Todo {
	for _, todo := range db.todos {
		if todo.ID == id {
			return todo
		}
	}
	return Todo{}
}

func (db DB) GetTodos() []Todo {
    log.Infof("Getting all todos: %v", db.todos)
    return db.todos
}

func (db DB) GetTodosByDone(isDone bool) []Todo {
    var todos []Todo
    for _,todo := range db.todos {
        if isDone == todo.Done {
            todos = append(todos, todo)
        }
    }
    return db.todos
}


func (db DB) GetRemainingTodoCount() string {
    var count int
    for _, todo := range db.todos {
        if !todo.Done {
            count++
        }
    }
    return strconv.Itoa(count)
}
