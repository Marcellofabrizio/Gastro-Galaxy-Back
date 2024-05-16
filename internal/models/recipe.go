package models

import (
	"database/sql"
)

type Recipe struct {
	id          int
	categoryId  int
	name        string
	url         string
	description string
}

type RecipeRepository struct {
	DB *sql.DB
}

func (r *RecipeRepository) Insert(name string, description string, url string, categoryId int) (int, error) {

	stmt := `INSERT INTO recipe (name, description, url, categoryId)
	VALUES(?,?,?,?)`

	result, err := r.DB.Exec(stmt, name, description, url, categoryId)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}
