package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Row interface {
	Scan(dest ...interface{}) error
}

func Initialize() error {
	tokens = make(map[string]*Token)

	var err error
	db, err = sql.Open("sqlite3", "dingo.db")
	if err != nil {
		return err
	}
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}
	checkBlogSettings()
	return nil
}

func checkBlogSettings() {
	SetSettingIfNotExists("theme", "default", "blog")
	SetSettingIfNotExists("title", "My Blog", "blog")
	SetSettingIfNotExists("description", "Awesome blog created by Dingo.", "blog")
}
