package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharykhin/todoapp/entity"
)

type TodoRepository struct {
	DB *sql.DB
}

func (tr TodoRepository) Get(limit int, offset int) ([]entity.Todo, error) {
	stmt, err := tr.DB.Prepare("SELECT * FROM todos LIMIT ? OFFSET ?")
	if err != nil {
		return nil, fmt.Errorf("could not make select statement:%v", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(limit, offset)
	var todos []entity.Todo
	for rows.Next() {
		var todo entity.Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed, &todo.Created)
		if err != nil {
			return nil, fmt.Errorf("error in scanning row to todo struct: %v", err)
		}
		//TODO: narrow case
		todos = append(todos, todo)
	}
	rows.Close()
	return todos, nil
}

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, fmt.Errorf("could not connect to sqlite database: %v", err)
	}

	return db, err
}
