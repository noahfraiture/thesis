package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := initDB()
	defer db.sqlite.Close()
	server(db)
}
