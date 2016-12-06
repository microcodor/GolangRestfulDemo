package main

import (
	"database/sql"
	"fmt"

	_ "../github.com/go-sql-driver/mysql"
)

var db *sql.DB

/*初始化数据库*/
func init() {
	//数据库操作
	db, _ = sql.Open("mysql", "root:xxxxxxx@/db_wordpress")
	db.Ping()
}

var currentId int

var todos Todos

func RepoFindTodo(id int) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateTodo(t Todo) Todo {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)
	return t
}

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
