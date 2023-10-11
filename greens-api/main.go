package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DBUser = "cameroncronheimer"
	DBPass = "camerondalton"
	DBName = "greensdb"
	DBHost = "localhost"
)

// Course struct to encode to JSON and send to the client
type Course struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Geometry string `json:"wkb_geometry"` // the polygon of the course
}

// Green struct
type Green struct {
	ID       string `json:"id"`           // primary key
	CID      string `json:"cid"`          // foreign key to course id (course.id)
	Geometry string `json:"wkb_geometry"` // the polygon of the green
	Centroid string `json:"centroid"`     // center of the green
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/courses", GetCourses).Methods("GET")
	r.HandleFunc("/greens", GetGreens).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // defer means this line will be executed once the function completes

	// Query the database
	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and print out the results
	for rows.Next() {
		var course Course                                            // create a new course struct for each row to scan into (encode to JSON)
		err := rows.Scan(&course.ID, &course.Name, &course.Geometry) // copy the columns from the row into the struct fields
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(course) // serialize the course struct to JSON and send it to the ResponseWriter
	}
}

func GetGreens(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // defer means this line will be executed once the function completes

	// Query the database
	rows, err := db.Query("SELECT * FROM greens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and print out the results
	for rows.Next() {
		var green Green
		err := rows.Scan(&green.ID, &green.CID, &green.Geometry, &green.Centroid)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(green)
	}
}

// ConnectDB connects to the database and returns a pointer to the database
func ConnectDB() (*sql.DB, error) {
	// Set up the connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", DBUser, DBPass, DBName, DBHost)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil

}
