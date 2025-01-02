package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Mountain struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Height    int    `json:"height"`
	LocalName string `json:"local_name"`
}

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PWD"),
		Net:    "tcp",
		Addr:   strings.Join([]string{os.Getenv("DB_HOST"), "3306"}, ":"),
		DBName: os.Getenv("DB_DATABASE"),
	}
	var err error
	// sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/dbname")
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/mountains", getMountains)
	router.POST("/mountains", addMountain)
	router.GET("/mountains/:id", getMountainByID)
	router.GET("/height/:h", getMountainsByHeight)
	router.Run("localhost:8080")

}

func getMountains(c *gin.Context) {
	// An mountains slice to hold data from returned rows.
	var mountains []Mountain
	var status int

	rows, err := db.Query("SELECT * FROM mountain")
	if err != nil {
		status = http.StatusNoContent
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var mtn Mountain
		if err := rows.Scan(&mtn.ID, &mtn.Name, &mtn.Height, &mtn.LocalName); err != nil {
			status = http.StatusInternalServerError
		}
		mountains = append(mountains, mtn)
	}
	// checking for an error here is the only way to find out that the results are incomplete,
	// i.e. if the query itself fails
	if err := rows.Err(); err != nil {
		status = http.StatusGatewayTimeout
	}
	c.IndentedJSON(status, mountains)
}

// mountainsByHeight queries for mountains that have the specified height.
func getMountainsByHeight(c *gin.Context) {
	// An mountains slice to hold data from returned rows.
	var mountains []Mountain
	height := c.Param("h")

	rows, err := db.Query("SELECT * FROM mountain WHERE height = ?", height)
	if err != nil {
		fmt.Printf("mountainsByHeight %q: %v\n", height, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var mtn Mountain
		if err := rows.Scan(&mtn.ID, &mtn.Name, &mtn.Height, &mtn.LocalName); err != nil {
			fmt.Printf("mountainsByHeight %q: %v\n", height, err)
		}
		mountains = append(mountains, mtn)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("mountainsByHeight %q: %v\n", height, err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "mountain(s) not found"})
	}
	c.IndentedJSON(http.StatusOK, mountains)
}

// getMountainByID queries for the mountain with the specified ID.
func getMountainByID(c *gin.Context) {
	// An mountain to hold data from the returned row.
	var mtn Mountain
	id := c.Param("id")

	row := db.QueryRow("SELECT * FROM mountain WHERE id = ?", id)
	if err := row.Scan(&mtn.ID, &mtn.Name, &mtn.Height, &mtn.LocalName); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "mountain not found"})
			fmt.Printf("mountainsById %s: no such mountain\n", id)
		}
		c.Status(http.StatusNotFound)
		fmt.Printf("mountainsById %s: %v\n", id, err)
	}
	c.IndentedJSON(http.StatusOK, mtn)
}

// addMountain adds the specified mountain to the database,
// returning the mountain ID of the new entry
func addMountain(c *gin.Context) {
	var newMountain Mountain

	// Call BindJSON to bind the received JSON to newMountain.
	if err := c.BindJSON(&newMountain); err != nil {
		return
	}
	result, err := db.Exec(
		"INSERT INTO mountain (name, height, local_name) VALUES (?, ?, ?)",
		newMountain.Name, newMountain.Height, newMountain.LocalName,
	)
	if err != nil {
		fmt.Printf("addMountain: %v\n", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "mountain not added"})
	}
	newMountain.ID = id
	c.IndentedJSON(http.StatusCreated, newMountain)
}
