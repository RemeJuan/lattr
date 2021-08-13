package domain

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTweetRepo_Get(t *testing.T) {
	var createdAt = time.Now()
	var modified = time.Now()

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := InitTweetRepository(db)

	var message = "message"
	var postTime = "2021-07-12 10:55:50 +0000 UTC"
	var userId = "001"
	var recordId int64 = 001

	expected := &Tweet{
		Id:        1,
		UserId:    userId,
		Message:   message,
		PostTime:  postTime,
		Status:    Pending,
		CreatedAt: createdAt,
		Modified:  modified,
	}

	t.Run("Success", func(t *testing.T) {
		const wantError = false

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, Pending, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, err := s.Get(recordId)

		if (err != nil) != wantError {
			t.Errorf("Get() error new = %v", err)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		const wantError = true
		expected := Tweet{}
		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, err := s.Get(recordId)

		if (err != nil) != wantError {
			t.Errorf("Get() error new = %v, wantErr %v", err, true)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})

	t.Run("Invalid Prepare", func(t *testing.T) {
		const wantError = true
		expected := Tweet{}
		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, Pending, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM wrongTable"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, err := s.Get(recordId)

		if (err != nil) != wantError {
			t.Errorf("Get() error new = %v, wantErr %v", err, true)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})
}

func TestTweetRepo_Create(t *testing.T) {
	var createdAt = time.Now()

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := InitTweetRepository(db)

	var message = ""
	var postTime = "2021-07-12 10:55:50 +0000 UTC"
	var userId = "001"
	var recordId int64 = 1

	expected := &Tweet{
		Id:        recordId,
		UserId:    userId,
		Message:   message,
		PostTime:  postTime,
		Status:    Pending,
		CreatedAt: createdAt,
	}

	request := &Tweet{
		UserId:    userId,
		Message:   message,
		PostTime:  postTime,
		CreatedAt: createdAt,
	}

	t.Run("Success", func(t *testing.T) {
		message = "message"
		const wantError = false

		expected.Message = message

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := sqlmock.NewResult(1, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnResult(sqlReturn)

		request.Message = message

		got, err := s.Create(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err.Message())
			t.Errorf("Create() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Create() = %v, want %v", got, expected)
		}
	})

	t.Run("Empty Message", func(t *testing.T) {
		message = ""
		const wantError = true

		expected.Message = message

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := errors.New("empty title")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnError(sqlReturn)

		request.Message = message

		got, err := s.Create(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err.Message())
			t.Errorf("Create() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Create() = %v, want %v", got, expected)
		}
	})

	t.Run("Invalid SQL query", func(t *testing.T) {
		const wantError = true

		sqlQuery := "INESRT INTO tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnError(sqlReturn)

		got, err := s.Create(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err.Message())
			t.Errorf("Create() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Create() = %v, want %v", got, expected)
		}
	})
}
