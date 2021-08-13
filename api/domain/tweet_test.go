package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var createdAt = time.Now()
var modified = time.Now()

func TestTweetRepo_Get(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := InitTweetRepository(db)

	t.Run("Success", func(t *testing.T) {
		const message = "message"
		const postTime = "2021-07-12 10:55:50 +0000 UTC"
		const userId = "001"
		const recordId = 001

		expected := &Tweet{
			Id:        1,
			UserId:    userId,
			Message:   message,
			PostTime:  postTime,
			Status:    Pending,
			CreatedAt: createdAt,
			Modified:  modified,
		}
		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, Pending, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, err := s.Get(recordId)

		if err != nil {
			t.Errorf("Get() error new = %v", err)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})

	//t.Run("Not Found", func(t *testing.T) {})
	//t.Run("Invalid Prepare", func(t *testing.T) {})
}
