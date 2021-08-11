package postgress

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func Connect() *gorm.DB {
	// open database
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	CheckError(err)

	defer db.Close()

	fmt.Println("Connected!")

	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
