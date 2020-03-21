package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "orm.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("DROP TABLE IF EXISTS User;")
	db.Exec("CREATE TABLE User(name TEXT);")

	result, err := db.Exec("INSERT into USER(name) values (?),(?)", "Foo", "Bar")
	if err == nil {
		// 返回受影响的行数
		n, _ := result.RowsAffected()
		fmt.Println("受影响的行数： ", n)
	}

	row := db.QueryRow("SELECT name from User LIMIT 1;")
	var name string
	if err = row.Scan(&name); err == nil {
		fmt.Println(name)
	}

}
