package domain

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/RemeJuan/lattr/utils/error_formats"
	"github.com/RemeJuan/lattr/utils/error_utils"
	_ "github.com/lib/pq"
)

var (
	TweetRepo TweetRepoInterface = &tweetRepo{}
)

var (
	queryGetTweet              = "SELECT Id, UserId, Message, PostTime, Status, CreatedAt, Modified FROM tweets WHERE id=$1;"
	queryInsertTweet           = "INSERT INTO tweets(UserId, Message, PostTime, Status, CreatedAt, Modified) VALUES($1, $2, $3, $4, $5, $6) RETURNING ID;"
	queryUpdateTweet           = "UPDATE tweets SET Message=$1, PostTime=$2 Status=$3 Modified=$4 WHERE id=$5;"
	queryGetAllTweets          = "SELECT * FROM tweets WHERE UserId=$1;"
	queryDeleteTweet           = "DELETE FROM tweets WHERE id=$1;"
	queryGetPendingTweets      = "SELECT * FROM tweets WHERE Status != 'Pending' AND PostTime <= now() order by PostTime asc LIMIT 1"
	queryGetLastScheduledTweet = "SELECT PostTime FROM tweets WHERE Status='Scheduled' ORDER by PostTime desc LIMIT 1"
)

type TweetRepoInterface interface {
	Initialize() *sql.DB
	Create(*Tweet) (*Tweet, error_utils.MessageErr)
	Get(int64) (*Tweet, error_utils.MessageErr)
	GetAll(string) ([]Tweet, error_utils.MessageErr)
	Update(*Tweet) (*Tweet, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	GetPending() ([]Tweet, error_utils.MessageErr)
	GetLast() (*Tweet, error_utils.MessageErr)
}

type tweetRepo struct {
	db *sql.DB
}

func InitTweetRepository(db *sql.DB) TweetRepoInterface {
	return &tweetRepo{
		db: db,
	}
}

func (tr *tweetRepo) Initialize() *sql.DB {
	var err error
	tr.db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	checkError(err)

	fmt.Println("Connected!")

	return tr.db
}

func (tr *tweetRepo) Create(tweet *Tweet) (*Tweet, error_utils.MessageErr) {
	var id int64
	stmt, err := tr.db.Prepare(queryInsertTweet)

	if err != nil {
		message := fmt.Sprintf("Error when trying to prepare all entries: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}
	defer stmt.Close()

	insertResult, createErr := stmt.Query(tweet.UserId, tweet.Message, tweet.PostTime, tweet.Status, tweet.CreatedAt, tweet.Modified)
	if createErr != nil {
		return nil, error_formats.ParseError(createErr)
	}

	insertResult.Next()
	inErr := insertResult.Scan(&id)

	if inErr != nil {
		message := fmt.Sprintf("error when trying to save data: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}

	tweet.Id = id
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
		return nil, error_formats.ParseError(getError)
	}

	return &tweet, nil
}

func (tr *tweetRepo) Update(tweet *Tweet) (*Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryUpdateTweet)

	if err != nil {
		message := fmt.Sprintf("error when trying to prepare update: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(tweet.Message, tweet.PostTime, tweet.Status, tweet.Modified, tweet.Id)
	if updateErr != nil {
		return nil, error_formats.ParseError(updateErr)
	}
	return tweet, nil
}

func (tr *tweetRepo) GetAll(userId string) ([]Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryGetAllTweets)

	if err != nil {
		return nil, error_utils.InternalServerError(fmt.Sprintf("Error when trying to prepare all entries: %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, error_formats.ParseError(err)
	}
	defer rows.Close()

	results := make([]Tweet, 0)

	for rows.Next() {
		var tweet Tweet
		if getError := rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Message, &tweet.PostTime, &tweet.Status, &tweet.CreatedAt, &tweet.Modified); getError != nil {
			message := fmt.Sprintf("Error when trying to get message: %s", getError.Error())
			return nil, error_utils.InternalServerError(message)
		}
		results = append(results, tweet)
	}
	if len(results) == 0 {
		return nil, error_utils.NotFoundError("no records found")
	}
	return results, nil
}

func (tr *tweetRepo) Delete(id int64) error_utils.MessageErr {
	stmt, err := tr.db.Prepare(queryDeleteTweet)
	if err != nil {
		return error_utils.InternalServerError(fmt.Sprintf("error when trying to delete record: %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		return error_utils.InternalServerError(fmt.Sprintf("error when trying to delete record %s", err.Error()))
	}
	return nil
}

func (tr *tweetRepo) GetPending() ([]Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryGetPendingTweets)

	if err != nil {
		return nil, error_utils.InternalServerError(fmt.Sprintf("Error when trying to prepare pending entries: %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, error_formats.ParseError(err)
	}
	defer rows.Close()

	results := make([]Tweet, 0)

	for rows.Next() {
		var tweet Tweet
		if getError := rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Message, &tweet.PostTime, &tweet.Status, &tweet.CreatedAt, &tweet.Modified); getError != nil {
			message := fmt.Sprintf("Error when trying to get message: %s", getError.Error())
			return nil, error_utils.InternalServerError(message)
		}
		results = append(results, tweet)
	}
	if len(results) == 0 {
		return nil, error_utils.NotFoundError("no records found")
	}
	return results, nil
}

func (tr *tweetRepo) GetLast() (*Tweet, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryGetLastScheduledTweet)

	if err != nil {
		message := fmt.Sprintf("Error retrieving record: %s", err)
		return nil, error_utils.InternalServerError(message)
	}

	defer stmt.Close()

	var tweet Tweet
	result := stmt.QueryRow()

	if getError := result.Scan(&tweet.PostTime); getError != nil {
		return nil, error_formats.ParseError(getError)
	}

	return &tweet, nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
