package domain

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/RemeJuan/lattr/utils/error_formats"
	"github.com/RemeJuan/lattr/utils/error_utils"
)

var (
	TweetRepo TweetRepoInterface = &tweetRepo{}
)

var (
	getTweetQuery = `SELECT Id, UserId, Message, PostTime, Status, CreatedAt, Modified FROM tweets WHERE id=?;`
)

type TweetRepoInterface interface {
	Initialize() *sql.DB
	Get(int64) (*Tweet, error_utils.MessageErr)
}

type tweetRepo struct {
	db *sql.DB
}

func (tr *tweetRepo) Initialize() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	checkError(err)

	defer db.Close()

	fmt.Println("Connected!")

	return db
}

func InitTweetRepository(db *sql.DB) TweetRepoInterface {
	return &tweetRepo{
		db: db,
	}
}

func (tr *tweetRepo) Get(id int64) (*Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(getTweetQuery)
	if err != nil {
		fmt.Println(err)
		message := fmt.Sprintf("Error retrieving record: %s", err)
		return nil, error_utils.InternalServerError(message)
	}

	defer stmt.Close()

	var tweet Tweet
	result := stmt.QueryRow(id)

	if getError := result.Scan(&tweet.Id, &tweet.UserId, &tweet.Message, &tweet.PostTime, &tweet.Status, &tweet.CreatedAt, &tweet.Modified); getError != nil {
		fmt.Println("this is the error man: ", getError)
		return nil, error_formats.ParseError(getError)
	}

	return &tweet, nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
