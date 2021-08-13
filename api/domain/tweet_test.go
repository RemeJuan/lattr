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

	const postTime = "2021-07-12 10:55:50 +0000 UTC"
	const userId = "001"
	const recordId int64 = 001
	var message = "message"

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

	const postTime = "2021-07-12 10:55:50 +0000 UTC"
	const userId = "001"
	const recordId int64 = 1
	var message = ""

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
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		message = "message"
		const wantError = false

		expected.Message = message

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := sqlmock.NewResult(1, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnResult(sqlReturn)

		request.Message = message

		got, err := s.Create(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Create() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Create() = %v, want %v", got, expected)
		}
	})

	t.Run("Empty Message", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		message = ""
		const wantError = true

		expected.Message = message

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := errors.New("empty title")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnError(sqlReturn)

		request.Message = message

		got, err := s.Create(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Create() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Create() = %v, want %v", got, expected)
		}
	})

	t.Run("Invalid SQL query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		sqlQuery := "INESRT INTO tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnError(sqlReturn)

		got, err := s.Create(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Create() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Create() = %v, want %v", got, expected)
		}
	})
}

func TestTweetRepo_Update(t *testing.T) {
	var modified = time.Now()

	const postTime = "2021-07-12 10:55:50 +0000 UTC"
	const recordId int64 = 001
	const status = Posted
	var message = "message"

	expected := &Tweet{
		Id:       recordId,
		Message:  message,
		PostTime: postTime,
		Status:   status,
		Modified: modified,
	}

	request := &Tweet{
		Id:       recordId,
		Message:  message,
		PostTime: postTime,
		Status:   status,
		Modified: modified,
	}

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = false

		const sqlQuery = "UPDATE tweets"
		sqlReturn := sqlmock.NewResult(0, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnResult(sqlReturn)

		got, err := s.Update(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Update() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Update() = %v, want %v", got, expected)
		}
	})

	t.Run("Invalid SQL Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("error in sql query statement")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, err := s.Update(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Update() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Update() = %v, want %v", got, expected)
		}
	})

	t.Run("Invalid Query Id", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		const sqlQuery = "UPDATE tweets"
		var sqlReturn = errors.New("invalid update id")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, err := s.Update(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Update() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Update() = %v, want %v", got, expected)
		}
	})

	t.Run("Empty Message", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		const sqlQuery = "UPDATE tweets"
		var sqlReturn = errors.New("please enter a valid title")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, 123).WillReturnError(sqlReturn)

		got, err := s.Update(request)

		if (err != nil) != wantError {
			fmt.Println("this is the error message: ", err)
			t.Errorf("Update() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Update() = %v, want %v", got, expected)
		}
	})

	t.Run("Update failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		const sqlQuery = "UPDATE tweets"
		sqlReturn := sqlmock.NewErrorResult(errors.New("update failed"))
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, 123).WillReturnResult(sqlReturn)

		got, err := s.Update(request)

		fmt.Println("this is the error message: ", err)
		if (err != nil) != wantError {
			t.Errorf("Update() error = %v, wantErr %v", err, wantError)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Update() = %v, want %v", got, expected)
		}
	})
}

func TestTweetRepo_GetAll(t *testing.T) {
	var modified = time.Now()
	var createdAt = time.Now()

	const postTime = "2021-07-12 10:55:50 +0000 UTC"
	const recordId int64 = 001
	const status = Posted
	const userId = "001"
	var message = "message"

	expected := []Tweet{
		{
			Id:        recordId,
			UserId:    userId,
			Message:   message,
			PostTime:  postTime,
			Status:    status,
			CreatedAt: createdAt,
			Modified:  modified,
		},
		{
			Id:        002,
			UserId:    userId,
			Message:   message,
			PostTime:  postTime,
			Status:    status,
			CreatedAt: createdAt,
			Modified:  modified,
		},
	}

	request := userId

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = false

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, status, createdAt, modified).AddRow(002, userId, message, postTime, status, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, err := s.GetAll(request)

		if (err != nil) != wantError {
			t.Errorf("Get() error new = %v", err)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})

	t.Run("Invalid SQL Syntax", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		const sqlQuery = "SELECTS (.+) FROM tweets"
		sqlResult := errors.New("error when trying to prepare all messages")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnError(sqlResult)

		got, err := s.GetAll(request)

		if (err != nil) != wantError {
			t.Errorf("Get() error new = %v", err)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})

	t.Run("No results", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const wantError = true

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, err := s.GetAll(request)

		if (err != nil) != wantError {
			t.Errorf("Get() error new = %v", err)
			return
		}

		if err == nil && !reflect.DeepEqual(got, expected) {
			t.Errorf("Get() = %v, want %v", got, expected)
		}
	})
}
