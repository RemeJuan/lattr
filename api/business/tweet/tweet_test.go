package tweet

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestScheduleTweetSuccess(t *testing.T) {
	const message = "message"
	const postTime = "2021-07-12 12:55:50 +0200"

	db, mock, err := sqlmock.New()
	gdb, gerr := gorm.Open("postgres", db)
	data := make(map[string]string)

	data["message"] = message
	data["time"] = postTime

	tweet := Tweet{
		Message:  message,
		UserId:   "001",
		PostTime: "2021-07-12 10:55:50 +0000 UTC",
		Status:   Pending,
	}

	if err != nil || gerr != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer gdb.Close()

	const sqlInsert = `INSERT INTO "tweets" ("created_at","updated_at","deleted_at","message","user_id","post_time","status") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "tweets"."id"`

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).WithArgs(AnyTime{}, AnyTime{}, tweet.DeletedAt, tweet.Message, tweet.UserId, tweet.PostTime, tweet.Status)

	_ = ScheduleTweet(gdb, data)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestScheduleTweetError(t *testing.T) {
	const message = "message"
	const postTime = "2021-07-12 12:55:50 0200"

	db, _, _ := sqlmock.New()
	gdb, _ := gorm.Open("postgres", db)
	data := make(map[string]string)

	data["message"] = message
	data["time"] = postTime

	shedErr := ScheduleTweet(gdb, data)

	if shedErr == nil {
		t.Log("Error should be nil")
		t.Fail()
	}
}

func TestBuildTweetSuccess(t *testing.T) {
	data := make(map[string]string)
	data["message"] = "hello"
	data["time"] = "2021-07-17 12:55:50 +0200"

	expected := Tweet{
		Message:  "hello",
		UserId:   "001",
		PostTime: "2021-07-17 10:55:50 +0000 UTC",
		Status:   Pending,
	}
	t2, err := BuildTweet(data)

	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
		return
	}

	if t2 != expected {
		t.Log("Output does not match expected", t2)
		t.Fail()
	}
}

func TestBuildTweetError(t *testing.T) {
	const message = "message"
	const postTime = "2021-07-12 12:55:50 0200"

	data := make(map[string]string)

	data["message"] = message
	data["time"] = postTime

	_, shedErr := BuildTweet(data)

	if shedErr == nil {
		t.Log("Error should be nil")
		t.Fail()
	}
}
