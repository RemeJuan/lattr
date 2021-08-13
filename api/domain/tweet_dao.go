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
	queryGetTweet    = "SELECT Id, UserId, Message, PostTime, Status, CreatedAt, Modified FROM tweets WHERE id=?;"
	queryInsertTweet = "INSERT INTO tweets(UserId, Message, PostTime, Status, CreatedAt) VALUES(?, ?, ?, ?, ?);"
	queryUpdateTweet = "UPDATE tweets SET Message=?, PostTime=? Status=? Modified=? WHERE id=?;"
)

type TweetRepoInterface interface {
	Initialize() *sql.DB
	Create(*Tweet) (*Tweet, error_utils.MessageErr)
	Get(int64) (*Tweet, error_utils.MessageErr)
	Update(*Tweet) (*Tweet, error_utils.MessageErr)
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

func (tr *tweetRepo) Create(tweet *Tweet) (*Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryInsertTweet)

	if err != nil {
		fmt.Println(err)
		message := fmt.Sprintf("Error when trying to prepare all messages: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}
	defer stmt.Close()

	insertResult, createErr := stmt.Exec(tweet.UserId, tweet.Message, tweet.PostTime, Pending, tweet.CreatedAt)
	if createErr != nil {
		return nil, error_formats.ParseError(createErr)
	}

	msgId, inErr := insertResult.LastInsertId()
	if inErr != nil {
		message := fmt.Sprintf("error when trying to save message: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}

	tweet.Id = msgId
	tweet.Status = Pending

	return tweet, nil
}

func (tr *tweetRepo) Get(id int64) (*Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryGetTweet)

	if err != nil {
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

func (tr *tweetRepo) Update(tweet *Tweet) (*Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryUpdateTweet)

	if err != nil {
		message := fmt.Sprintf("error when trying to prepare user to update: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(tweet.Message, tweet.PostTime, tweet.Status, tweet.Modified, tweet.Id)
	if updateErr != nil {
		return nil, error_formats.ParseError(updateErr)
	}
	return tweet, nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
