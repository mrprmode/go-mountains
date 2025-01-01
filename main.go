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
	// if the query itself fails
	if err := rows.Err(); err != nil {
		status = http.StatusGatewayTimeout
	}

	c.IndentedJSON(status, mountains)
}
