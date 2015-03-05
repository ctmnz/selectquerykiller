package main

import (
//	"flag"
	"os"
	"log"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)


func killquery(db *sql.DB, qid int,qinfo sql.NullString) {
	log.Println("Killing queryi: ", qinfo)
	err := db.QueryRow("kill ?", qid)
	if err != nil {
		log.Fatal(err)
	}
}




func main() {

	if (len(os.Args) < 2 ) {
		fmt.Println("Usage: ", os.Args[0] ," <SECONDS>")
		os.Exit(2)
	}

	waitseconds := os.Args[1]

	 waits, errint := strconv.Atoi(waitseconds)
	if errint != nil {
		fmt.Println("Illegal argument")
		os.Exit(2)
	}


	db, err := sql.Open("mysql", "root:@/test")
	if err != nil {
		fmt.Println("Error!: ", err)
	}

	rows, _ := db.Query("SELECT ID,TIME,INFO FROM information_schema.processlist WHERE INFO like \"SELECT %\";")
	defer rows.Close()

	for rows.Next() {
		var id int
		var time int
		var info sql.NullString

// Id    | User  | Host             | db       | Command | Time | State | Info          
		if err := rows.Scan( &id, &time, &info); err != nil {
			log.Fatal(err)
		}
		if (time > waits) {
			fmt.Printf("id: %d | time: %d | killing query: %s \n",id , time, info)
			killquery(db, id, info)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
