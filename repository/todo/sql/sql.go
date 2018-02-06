package sql

import (
	"database/sql"
	"fmt"

	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharykhin/todoapp/entity"
	"github.com/sharykhin/todoapp/request"
	"github.com/sharykhin/todoapp/utils"
)

type TodoRepository struct {
	DB *sql.DB
}

func (tr TodoRepository) Get(limit, offset string) ([]entity.Todo, error) {
	stmt, err := tr.DB.Prepare("SELECT * FROM todos LIMIT ? OFFSET ?")
	if err != nil {
		return nil, fmt.Errorf("could not make select statement:%v", err)
	}
	defer stmt.Close()

	var todos []entity.Todo
	rows, err := stmt.Query(limit, offset)
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

func (tr TodoRepository) Count() (*int, error) {
	var count string
	stmt, err := tr.DB.Prepare("SELECT COUNT(id) as `total` FROM todos")
	if err != nil {
		return nil, fmt.Errorf("could not make select statement:%v", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("could not make scan: %v", err)
	}
	i, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("could not convert string to number: %v", err)
	}
	return &i, nil
}

func (tr TodoRepository) Create(rt request.Todo) (*entity.Todo, error) {
	stmt, err := tr.DB.Prepare("INSERT INTO todos(title, description,completed,created) values(?,?,?,?)")
	if err != nil {
		return nil, fmt.Errorf("could not make insert statement:%v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(rt.Title, rt.Description, rt.Completed, rt.Created)
	if err != nil {
		return nil, fmt.Errorf("could not apply values: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not get last insert id: %v", err)
	}
	t := entity.Todo{
		Id:          int(id),
		Title:       rt.Title,
		Description: rt.Description,
		Completed:   rt.Completed,
		Created:     utils.JSONTime(rt.Created),
	}
	return &t, nil
}
