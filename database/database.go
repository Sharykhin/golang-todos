package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sharykhin/golang-todos/entity"
	"github.com/Sharykhin/golang-todos/utils"
	_ "github.com/mattn/go-sqlite3" // dependency of sqlite3
)

var db *sql.DB

// Storage variable is a link to a private struct that actually manages everything regarding todos
var Storage storage

// storage struct is responsible for managing all todos in a database
type storage struct{}

// Create calls a package method for creating a new item
func (s storage) Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
	res, err := db.ExecContext(ctx, "INSERT INTO todos(title, description, completed) VALUES (?, ?, ?)", rt.Title, rt.Description, rt.Completed)
	if err != nil {
		return nil, fmt.Errorf("could not make insert statement: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not get last insert id: %v", err)
	}
	t := entity.Todo{
		ID:          id,
		Title:       rt.Title,
		Description: rt.Description,
		Completed:   rt.Completed,
		Created:     utils.JSONTime(time.Now().UTC()),
	}
	return &t, nil
}

func (s storage) Get(ctx context.Context, limit, offset int) ([]entity.Todo, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM todos LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not make select statement: %v", err)
	}
	defer rows.Close()

	var todos []entity.Todo
	for rows.Next() {
		var todo entity.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.Created)
		if err != nil {
			return nil, fmt.Errorf("error in scanning row to todo struct: %v", err)
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

func (s storage) Count(ctx context.Context) (int, error) {
	var count int
	row := db.QueryRowContext(ctx, "SELECT COUNT(id) AS `total` FROM todos")
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not make scan: %v", err)
	}
	return count, nil
}

// DB returns current instance of database
func DB() *sql.DB {
	return db
}

func init() {
	var err error
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "./foo.db"
	}
	db, err = sql.Open("sqlite3", dbSource)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}
}

// Get returns limited todos
//func Get(ctx context.Context, limit, offset int) ([]entity.Todo, error) {
//	rows, err := db.QueryContext(ctx, "SELECT * FROM todos LIMIT ? OFFSET ?", limit, offset)
//	if err != nil {
//		return nil, fmt.Errorf("could not make select statement: %v", err)
//	}
//	defer rows.Close()
//
//	var todos []entity.Todo
//	for rows.Next() {
//		var todo entity.Todo
//		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.Created)
//		if err != nil {
//			return nil, fmt.Errorf("error in scanning row to todo struct: %v", err)
//		}
//		todos = append(todos, todo)
//	}
//	return todos, rows.Err()
//}
//
//// Count returns all count of todos
//func Count(ctx context.Context) (int, error) {
//	var count int
//	row := db.QueryRowContext(ctx, "SELECT COUNT(id) AS `total` FROM todos")
//	err := row.Scan(&count)
//	if err != nil {
//		return 0, fmt.Errorf("could not make scan: %v", err)
//	}
//	return count, nil
//}
//
//// Create creates new todo and returns new item
//func Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
//	res, err := db.ExecContext(ctx, "INSERT INTO todos(title, description, completed) VALUES (?, ?, ?)", rt.Title, rt.Description, rt.Completed)
//	if err != nil {
//		return nil, fmt.Errorf("could not make insert statement: %v", err)
//	}
//	id, err := res.LastInsertId()
//	if err != nil {
//		return nil, fmt.Errorf("could not get last insert id: %v", err)
//	}
//	t := entity.Todo{
//		ID:          id,
//		Title:       rt.Title,
//		Description: rt.Description,
//		Completed:   rt.Completed,
//		Created:     utils.JSONTime(time.Now().UTC()),
//	}
//	return &t, nil
//}
