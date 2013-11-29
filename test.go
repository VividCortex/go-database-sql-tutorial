package main

import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/VividCortex/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root@tcp(127.0.0.1:12830)/customers?charset=utf8")
	var res string
	err := db.QueryRow("call template.foo").Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
