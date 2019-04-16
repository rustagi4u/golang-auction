package shared

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // _ needed
)

const (
	DatabaseHost = "localhost"
	// DatabaseUser is the database user
	DatabaseUser = "postgres"
	// DatabasePassword is the password used
	DatabasePassword = "restgomuxpq"
	// DatabaseName ‘quest’ is the name of our database
	DatabaseName = "postgres"
)

// DB is exported db
type DB struct {
	*sql.DB
}

// NewOpen exported for creating a new Postgresql instance
// func Connect() (DB, error) {
func NewOpen() (DB, error) {
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DatabaseHost, DatabaseUser, DatabasePassword, DatabaseName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}

	return DB{db}, err
}
