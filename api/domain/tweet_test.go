package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTweetRepo_Initialize(t *testing.T) {
	t.Skipf("To DO")
}

const layout = "2021-07-12 10:55:50 +0000"

func TestTweetRepo_Create(t *testing.T) {
	var createdAt = time.Now().Local()
	var modified = time.Now().Local()

	const userId = "001"
	const recordId int64 = 1
	var message = ""
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	request := &Tweet{
		UserId:    userId,
		Message:   message,
		PostTime:  postTime,
		Status:    Pending,
		CreatedAt: createdAt,
		Modified:  modified,
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
			Modified:  modified,
		}

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := sqlmock.NewRows([]string{"Id"}).AddRow(recordId)
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId, message, postTime, Pending, createdAt, modified).WillReturnRows(sqlReturn)

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

		const expected = "error when trying to save data: empty title"

		message = ""

		sqlQuery := "INSERT INTO tweets"
		sqlReturn := errors.New("empty title")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId, message, postTime, Pending, createdAt, modified).WillReturnError(sqlReturn)

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

		const expected = "Error when trying to prepare all entries: invalid sql query"

		const sqlQuery = "INSERT INTO tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		got, crErr := s.Create(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, crErr.Message())
	})
}

func TestTweetRepo_Get(t *testing.T) {
	var createdAt = time.Now().Local()
	var modified = time.Now().Local()

	const userId = "001"
	const recordId int64 = 001
	var message = "message"
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, Pending, createdAt, modified)

		expected := &Tweet{
			Id:        1,
			UserId:    userId,
			Message:   message,
			PostTime:  postTime,
			Status:    Pending,
			CreatedAt: createdAt,
			Modified:  modified,
		}

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, getErr := s.Get(recordId)

		assert.Nil(t, getErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "no record matching given id"

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(recordId).WillReturnRows(rows)

		got, gotErr := s.Get(recordId)

		assert.Nil(t, got)
		assert.Equal(t, expected, gotErr.Message())
	})

	t.Run("Invalid Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "Error retrieving record: invalid sql query"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		got, gotErr := s.Get(recordId)

		assert.Nil(t, got)
		assert.Equal(t, expected, gotErr.Message())
	})
}

func TestTweetRepo_Update(t *testing.T) {
	var modified = time.Now().Local()

	const recordId int64 = 001
	const status = Posted
	var message = "message"
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

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

		expected := &Tweet{
			Id:       recordId,
			Message:  message,
			PostTime: postTime,
			Status:   status,
			Modified: modified,
		}

		const sqlQuery = "UPDATE tweets"
		sqlReturn := sqlmock.NewResult(0, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnResult(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, upErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Invalid SQL Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "error when trying to prepare update: error in sql query statement"

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("error in sql query statement")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})

	t.Run("Invalid Query Id", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "error when trying to save data: invalid update id"

		const sqlQuery = "UPDATE tweets"
		var sqlReturn = errors.New("invalid update id")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})

	t.Run("Empty Message", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "error when trying to save data: please enter a valid title"

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("please enter a valid title")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})

	t.Run("Update failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "error when trying to save data: update failed"

		const sqlQuery = "UPDATE tweets"
		sqlReturn := errors.New("update failed")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(message, postTime, status, modified, recordId).WillReturnError(sqlReturn)

		got, upErr := s.Update(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, upErr.Message())
	})
}

func TestTweetRepo_GetAll(t *testing.T) {
	var modified = time.Now().Local()
	var createdAt = time.Now().Local()

	const recordId int64 = 001
	const status = Posted
	const userId = "001"
	var message = "message"
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	request := userId

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

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

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, status, createdAt, modified).AddRow(002, userId, message, postTime, status, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, gaErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Invalid SQL Syntax", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "Error when trying to prepare all entries: invalid syntax"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlResult := errors.New("invalid syntax")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlResult)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("Invalid Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "error when trying to save data: invalid query"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlResult := errors.New("invalid query")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnError(sqlResult)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("Invalid response.incorrect row", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		expected := "Error when trying to get message: sql: Scan error on column index 5, name \"CreatedAt\": unsupported Scan, storing driver.Value type string into type *time.Time"

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, status, createdAt, modified).AddRow(002, userId, message, postTime, status, "createdAt", modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("No results", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "no records found"

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(userId).WillReturnRows(rows)

		got, gaErr := s.GetAll(request)

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})
}

func TestTweetRepo_Delete(t *testing.T) {
	const recordId int64 = 1

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const sqlQuery = "DELETE FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(recordId).WillReturnResult(sqlmock.NewResult(0, 1))

		delErr := s.Delete(recordId)

		assert.Nil(t, delErr)
	})

	t.Run("Invalid Id/Not Found Id", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		expected := "error when trying to delete record id not found or invalid"

		const sqlQuery = "DELETE FROM tweets"
		sqlResult := errors.New("id not found or invalid")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(recordId).WillReturnError(sqlResult)

		delErr := s.Delete(recordId)

		assert.Equal(t, expected, delErr.Message())
	})

	t.Run("Invalid SQL query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		expected := "error when trying to delete record: id not found or invalid"

		const sqlQuery = "DELETE FROM tweets"
		sqlResult := errors.New("id not found or invalid")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlResult)

		delErr := s.Delete(recordId)

		assert.Equal(t, expected, delErr.Message())
	})
}

func TestTweetRepo_GetPending(t *testing.T) {
	var modified = time.Now().Local()
	var createdAt = time.Now().Local()

	const recordId int64 = 001
	const status = Posted
	const userId = "001"
	var message = "message"
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

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

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, status, createdAt, modified).AddRow(002, userId, message, postTime, status, createdAt, modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(rows)

		got, gaErr := s.GetPending()

		assert.Nil(t, gaErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Invalid SQL Syntax", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "Error when trying to prepare pending entries: invalid syntax"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlResult := errors.New("invalid syntax")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlResult)

		got, gaErr := s.GetPending()

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("Invalid Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "error when trying to save data: invalid query"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlResult := errors.New("invalid query")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnError(sqlResult)

		got, gaErr := s.GetPending()

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("Invalid response.incorrect row", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		expected := "Error when trying to get message: sql: Scan error on column index 5, name \"CreatedAt\": unsupported Scan, storing driver.Value type string into type *time.Time"

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"}).AddRow(recordId, userId, message, postTime, status, createdAt, modified).AddRow(002, userId, message, postTime, status, "createdAt", modified)

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(rows)

		got, gaErr := s.GetPending()

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})

	t.Run("No results", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "no records found"

		rows := sqlmock.NewRows([]string{"Id", "UserId", "Message", "PostTime", "Status", "CreatedAt", "Modified"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(rows)

		got, gaErr := s.GetPending()

		assert.Nil(t, got)
		assert.Equal(t, expected, gaErr.Message())
	})
}

func TestTweetRepo_GetLast(t *testing.T) {
	postTime, _ := time.Parse(layout, "2021-07-12 10:55:50 +0000")

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		rows := sqlmock.NewRows([]string{"PostTime"}).AddRow(postTime)

		expected := &Tweet{
			PostTime: postTime,
		}

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(rows)

		got, getErr := s.GetLast()

		assert.Nil(t, getErr)
		assert.Equal(t, expected, got)
	})

	t.Run("Not Found", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "no record matching given id"

		rows := sqlmock.NewRows([]string{"PostTime"})

		const sqlQuery = "SELECT (.+) FROM tweets"
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(rows)

		got, gotErr := s.GetLast()

		assert.Nil(t, got)
		assert.Equal(t, expected, gotErr.Message())
	})

	t.Run("Invalid Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTweetRepository(db)

		const expected = "Error retrieving record: invalid sql query"

		const sqlQuery = "SELECT (.+) FROM tweets"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		got, gotErr := s.GetLast()

		assert.Nil(t, got)
		assert.Equal(t, expected, gotErr.Message())
	})
}
