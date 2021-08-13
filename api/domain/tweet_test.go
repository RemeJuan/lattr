package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTweetRepo_Create(t *testing.T) {
	var createdAt = time.Now()

	const postTime = "2021-07-12 10:55:50 +0000 UTC"
	const userId = "001"
	const recordId int64 = 1
	var message = ""

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

		expected := &Tweet{
			Id:        recordId,
			UserId:    userId,
			Message:   message,
			PostTime:  postTime,
			Status:    Pending,
			CreatedAt: createdAt,
		}

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := sqlmock.NewResult(1, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnResult(sqlReturn)

		request.Message = message

		got, crErr := s.Create(request)

		assert.Nil(t, crErr)
		assert.Equal(t, got, expected)
	})

	t.Run("Empty Message", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		message = ""

		expected := "error when trying to save message: empty title"

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := errors.New("empty title")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnError(sqlReturn)

		request.Message = message

		got, createError := s.Create(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, createError.Message())
	})

	t.Run("Invalid SQL query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		expected := "error when trying to save message: invalid sql query"

		const sqlQuery = "INSERT INTO tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(userId, message, postTime, Pending, createdAt).WillReturnError(sqlReturn)

		got, crErr := s.Create(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, crErr.Message())
	})
}

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

	t.Run("Success", func(t *testing.T) {
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

		got, getErr := s.Get(recordId)

		assert.Nil(t, getErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Not Found", func(t *testing.T) {
		expected := "no record matching given id"
		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, gotErr := s.Get(recordId)

		assert.Nil(t, got)
		assert.Equal(t, expected, gotErr.Message())
	})

	t.Run("Invalid Prepare", func(t *testing.T) {
		const expected = "error when trying to save message: invalid sql query"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnError(sqlReturn)

		got, gotErr := s.Get(recordId)

		assert.Nil(t, got)
		assert.Equal(t, expected, gotErr.Message())
	})
}

func TestTweetRepo_Update(t *testing.T) {
	var modified = time.Now()

	const postTime = "2021-07-12 10:55:50 +0000 UTC"
	const recordId int64 = 001
	const status = Posted
	var message = "message"

	request := &Tweet{
		Id:       recordId,
		Message:  message,
		PostTime: postTime,
		Status:   status,
		Modified: modified,
	}

	t.Run("Success", func(t *testing.T) {
		expected := &Tweet{
			Id:       recordId,
			Message:  message,
			PostTime: postTime,
			Status:   status,
			Modified: modified,
		}

		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "UPDATE tweets"
		sqlReturn := sqlmock.NewResult(0, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnResult(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, upErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Invalid SQL Query", func(t *testing.T) {
		const expected = "error when trying to save data: error in sql query statement"
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("error in sql query statement")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})

	t.Run("Invalid Query Id", func(t *testing.T) {
		const expected = "error when trying to save data: invalid update id"
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "UPDATE tweets"
		var sqlReturn = errors.New("invalid update id")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})

	t.Run("Empty Message", func(t *testing.T) {
		const expected = "error when trying to save data: please enter a valid title"

		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("please enter a valid title")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})

	t.Run("Update failed", func(t *testing.T) {
		const expected = "error when trying to save data: update failed"

		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("update failed")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
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

	request := userId

	t.Run("Success", func(t *testing.T) {
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

		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, status, createdAt, modified).AddRow(002, userId, message, postTime, status, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, gaErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Invalid SQL Syntax", func(t *testing.T) {
		const expected = "error when trying to save data: error when trying to prepare all messages"
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlResult := errors.New("error when trying to prepare all messages")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnError(sqlResult)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("No results", func(t *testing.T) {
		const expected = "no records found"

		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})
}
