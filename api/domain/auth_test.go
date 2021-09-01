package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const mockToken = "mock-test-token"
const mockTokenName = "token-name"
const mockTokenId = 1

var ct = time.Now().Local()
var mt = time.Now().Local()
var exp = time.Now().Local()
var sc = []string{"token:create"}

func TestTokenRepo_Initialize(t *testing.T) {
	t.Skipf("To DO")
}

func TestTokenRepo_Create(t *testing.T) {
	request := &Token{
		Id:        mockTokenId,
		Name:      mockTokenName,
		Token:     mockToken,
		Scopes:    sc,
		ExpiresAt: exp,
		CreatedAt: ct,
		Modified:  mt,
	}

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		expected := &Token{
			Id:        mockTokenId,
			Name:      mockTokenName,
			Token:     mockToken,
			Scopes:    sc,
			ExpiresAt: exp,
			CreatedAt: ct,
			Modified:  mt,
		}
		const sqlQuery = "INSERT INTO tokens"
		sqlReturn := sqlmock.NewRows([]string{"token"}).AddRow(mockToken)
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(mockTokenName, mockToken, pq.Array(sc), exp, ct, mt).WillReturnRows(sqlReturn)

		result, crErr := s.Create(request)

		assert.Nil(t, crErr)
		assert.Equal(t, result, expected)
	})

	t.Run("Empty response", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to save data: empty name"
		const sqlQuery = "INSERT INTO tokens"
		sqlReturn := errors.New("empty name")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(mockTokenName, mockToken, pq.Array(sc), exp, ct, mt).WillReturnError(sqlReturn)

		result, crErr := s.Create(request)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("Invalid SQL Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "Error when trying to prepare all entries: invalid sql query"
		const sqlQuery = "INSERT INTO tokens"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		result, crErr := s.Create(request)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})
}

func TestTokenRepo_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		expected := &Token{
			Id:        mockTokenId,
			Name:      mockTokenName,
			Token:     mockToken,
			Scopes:    sc,
			ExpiresAt: exp,
			CreatedAt: ct,
			Modified:  mt,
		}
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := sqlmock.NewRows([]string{"id", "name", "token", "scopes", "expiresAt", "createdAt", "modified"}).AddRow(mockTokenId, mockTokenName, mockToken, pq.Array(sc), exp, ct, mt)
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(mockTokenId).WillReturnRows(sqlReturn)

		result, crErr := s.Get(mockTokenId)

		assert.Nil(t, crErr)
		assert.Equal(t, result, expected)
	})

	t.Run("Not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "no record matching given id"
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := sqlmock.NewRows([]string{"id", "name", "token", "createdAt", "modified"})
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WithArgs(mockTokenId).WillReturnRows(sqlReturn)

		result, crErr := s.Get(mockTokenId)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("Invalid prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "Error retrieving record: invalid sql query"
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := errors.New("invalid sql query")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		result, crErr := s.Get(mockTokenId)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})
}

func TestTokenRepo_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		expected := []Token{
			{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: ct,
				Modified:  mt,
			},
			{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: ct,
				Modified:  mt,
			},
		}
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := sqlmock.NewRows([]string{"id", "name", "token", "createdAt", "modified"}).AddRow(mockTokenId, mockTokenName, mockToken, ct, mt).AddRow(mockTokenId, mockTokenName, mockToken, ct, mt)
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(sqlReturn)

		result, crErr := s.List()

		assert.Nil(t, crErr)
		assert.Equal(t, expected, result)
	})

	t.Run("Invalid Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "Error when trying to prepare all entries: invalid syntax"
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := errors.New("invalid syntax")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		result, crErr := s.List()

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("Invalid Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to save data: invalid query"
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := errors.New("invalid query")
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnError(sqlReturn)

		result, crErr := s.List()

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("Invalid response incorrect row", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "Error when trying to get message: sql: Scan error on column index 3, name \"createdAt\": unsupported Scan, storing driver.Value type string into type *time.Time"
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := sqlmock.NewRows([]string{"id", "name", "token", "createdAt", "modified"}).AddRow(mockTokenId, mockTokenName, mockToken, "CreatedAt", mt)
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(sqlReturn)

		result, crErr := s.List()

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("No results", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "no records found"
		const sqlQuery = "SELECT (.+) FROM tokens"
		sqlReturn := sqlmock.NewRows([]string{"id", "name", "token", "createdAt", "modified"})
		mock.ExpectPrepare(sqlQuery).ExpectQuery().WillReturnRows(sqlReturn)

		result, crErr := s.List()

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})
}

func TestTokenRepo_Reset(t *testing.T) {
	request := &Token{
		Id:        mockTokenId,
		Token:     mockToken,
		Scopes:    sc,
		ExpiresAt: exp,
		Modified:  mt,
	}

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		expected := &Token{
			Id:        mockTokenId,
			Token:     mockToken,
			Scopes:    sc,
			ExpiresAt: exp,
			Modified:  mt,
		}
		const sqlQuery = "UPDATE tokens"
		sqlReturn := sqlmock.NewResult(0, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(mockTokenId, mockToken, exp, mt).WillReturnResult(sqlReturn)

		result, crErr := s.Reset(request)

		assert.Nil(t, crErr)
		assert.Equal(t, expected, result)
	})

	t.Run("Invalid SQL Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to prepare update: error in sql query statement"

		const sqlQuery = "UPDATE tokens"
		sqlReturn := errors.New("error in sql query statement")
		mock.ExpectPrepare(sqlQuery).WillReturnError(sqlReturn)

		result, crErr := s.Reset(request)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("Invalid Query Id", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to save data: invalid update id"

		const sqlQuery = "UPDATE tokens"
		sqlReturn := errors.New("invalid update id")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(mockTokenId, mockToken, exp, mt).WillReturnError(sqlReturn)

		result, crErr := s.Reset(request)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})

	t.Run("Update failed", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to save data: update failed"
		const sqlQuery = "UPDATE tokens"
		sqlReturn := errors.New("update failed")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(mockTokenId, mockToken, exp, mt).WillReturnError(sqlReturn)

		result, crErr := s.Reset(request)

		assert.Nil(t, result)
		assert.Equal(t, expected, crErr.Message())
	})
}

func TestTokenRepo_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const sqlQuery = "DELETE FROM tokens"
		sqlReturn := sqlmock.NewResult(0, 1)
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(mockTokenId).WillReturnResult(sqlReturn)

		delErr := s.Delete(mockTokenId)

		assert.Nil(t, delErr)
	})

	t.Run("Invalid Id/Not Found Id", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to delete record id not found or invalid"
		const sqlQuery = "DELETE FROM tokens"
		sqlReturn := errors.New("id not found or invalid")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(mockTokenId).WillReturnError(sqlReturn)

		delErr := s.Delete(mockTokenId)

		assert.Equal(t, expected, delErr.Message())
	})

	t.Run("Invalid SQL query", func(t *testing.T) {
		db, mock, err := sqlmock.New()

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		s := InitTokenRepository(db)

		const expected = "error when trying to delete record invalid query"
		const sqlQuery = "DELETE FROM tokens"
		sqlReturn := errors.New("invalid query")
		mock.ExpectPrepare(sqlQuery).ExpectExec().WithArgs(mockTokenId).WillReturnError(sqlReturn)

		delErr := s.Delete(mockTokenId)

		assert.Equal(t, expected, delErr.Message())
	})
}
