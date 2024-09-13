package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync"
	"time"
)

var db *sql.DB
var wg sync.WaitGroup

func main() {

	// 1. provision db in makefile and connect
	var err error
	db, err = sql.Open("postgres", "user=postgres password=donotshare dbname=test sslmode=disable")
	// 2. handle error
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Create the 'users' table if it doesn't exist. This ideally should be in the make file
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100))`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createUser)

	log.Println("server started")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		rows, err := db.Query("SELECT name FROM users")
		// 4. Handle error and return with the write status code
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get users: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			rows.Scan(&name)
			fmt.Fprintf(w, "User: %s\n", name)
		}
	}()

	wg.Wait()

}

func createUser(w http.ResponseWriter, r *http.Request) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		time.Sleep(5 * time.Second) // Simulate a long repository operation

		username := r.URL.Query().Get("name")
		// 5. Handle error if username is not present in query
		if username == "" {
			http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
			return
		}

		// 6. re-write db query to guard against SQL injection attacks
		_, err := db.Exec("INSERT INTO users (name) VALUES ($1)", username)

		if err != nil {
			// 7. Handle error with the right status code
			http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User %s created successfully", username)
	}()

	wg.Wait()
}
