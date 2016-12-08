package main

import (
	"database/sql"
	"fmt"
	"time"

	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

/*初始化数据库*/
func init() {
	//数据库操作
	db, _ = sql.Open("mysql", "root:jinchun123#@/db_wordpress")
	db.Ping()

	//获取时间戳

	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	fmt.Println("time:", tm.Format("2006-01-02"))

	logger, err := seelog.LoggerFromConfigAsFile("seelog.xml")

	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)

	seelog.Error("seelog error")
	seelog.Info("seelog info")
	seelog.Debug("seelog debug")

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
