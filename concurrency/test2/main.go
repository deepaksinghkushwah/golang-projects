package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var wg sync.WaitGroup

func main() {
	start := time.Now()
	wg.Add(1)
	go func() {
		db, _ = sql.Open("sqlite3", "./contest.db")

		defer db.Close()

		/*stmt, err := db.Prepare("INSERT INTO `user`(username, email, password, fullname) VALUES(?, ?, ?, ?)")

		if err != nil {
			log.Fatalln(err)
		}
		defer stmt.Close()

		for i := 711; i <= 1000; i++ {
			str := "test" + strconv.Itoa(i)
			iusername := str
			ipassword := str
			iemail := str + "@localhost.com"
			ifullname := str
			_, err := stmt.Exec(iusername, iemail, ipassword, ifullname)
			if err != nil {
				log.Fatalln(err)
			}

		}*/

		rows, err := db.Query("SELECT id, username, password, fullname, email FROM `user`")

		if err != nil {
			if err == sql.ErrNoRows {
				log.Fatalln("No records found")
			} else {
				log.Fatalln(err)
			}
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var username, password, fullname, email string
			err = rows.Scan(&id, &username, &password, &fullname, &email)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, username, password, fullname, email)
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Execution took: ", time.Since(start))
}
