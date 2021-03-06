package scritps

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func mysqlTest() {
	fmt.Println("Hello, world!")
	// Connect to database
	db, err := sql.Open("mysql", "nathaniel:nathaniel@/NatTest")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping database
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping successful!")

	// Query for something
	resultRows, err := db.Query("SELECT name FROM Example WHERE age>=50")
	if err != nil {
		log.Fatal(err)
	}
	defer resultRows.Close()

	// Display the resulting rows from the query
	fmt.Println("\nQuery results:")
	for resultRows.Next() {
		var name string
		if err = resultRows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)
	}

	// Error-check after query
	if err := resultRows.Err(); err != nil {
		log.Fatal(err)
	}

}
