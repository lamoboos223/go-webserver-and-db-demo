package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var apikeys []string

type Message struct {
	Text string `json:"text"`
}

func connectDB() (*sql.DB, error) {
	// Connection string
	connStr := "host=localhost port=5432 username=postgres password=mysecretpassword dbname=postgres sslmode=disable"

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Check if the connection is successful
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	return db, nil
}
func fetchApikeys(db *sql.DB, ch chan []string) {

	apikeys := make([]string, 0, 10)
	rows, err := db.Query("SELECT apikey FROM authentication_keys")
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var apikey string
		if err := rows.Scan(&apikey); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("apikey: %s\n", apikey)
		apikeys = append(apikeys, apikey)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}
	ch <- apikeys // Send apikeys through the channel
	close(ch)
}
func stringInSlice(target string, list []string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Parse JSON request body
	var message Message
	apikey := r.Header.Get("apikey")

	if err := apikey; err == "" {
		errorMessage := "No Api key given"
		http.Error(w, errorMessage, http.StatusUnauthorized)
		return
	} else {
		fmt.Println("apikey is: ", apikey)
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if stringInSlice(apikey, apikeys) {
		// Validate and process the message if needed
		// For demonstration, let's just echo back the message
		response := Message{Text: "Hello, " + message.Text}

		// Encode response as JSON and send it
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		response := Message{Text: "You have to register your api key"}
		// manipulate status code
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func main() {

	// Connect to PostgreSQL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Fetch users
	// Receive apikeys from the channel
	// Use a channel to receive apikeys
	apikeysChan := make(chan []string)
	go fetchApikeys(db, apikeysChan)
	apikeys = <-apikeysChan
	db.Close()

	// run the wevserver in Goroutines mode to enhance performance
	go http.HandleFunc("/hello", helloWorldHandler)

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// log.Fatal it's like system.exit(2)
		log.Fatal(err)
	}
}
