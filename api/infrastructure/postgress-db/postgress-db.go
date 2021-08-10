package postgress_db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type pgConnection struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func Connect() *gorm.DB {
	connection := connectionDetails()
	host := connection.host
	port := connection.port
	user := connection.user
	password := connection.password
	dbname := connection.dbname

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	// open database
	db, err := gorm.Open("postgres", psqlconn)
	CheckError(err)

	// Close the Connection in the calling function

	fmt.Println("Connected!")

	return db
}

func connectionDetails() pgConnection {
	port, _ := strconv.ParseInt(os.Getenv("PG_PORT"), 10, 32)

	return pgConnection{
		host:     os.Getenv("PG_HOST"),
		port:     int(port),
		user:     os.Getenv("PG_USER"),
		password: os.Getenv("PG_PASSWORD"),
		dbname:   os.Getenv("PG_DB_NAME"),
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
