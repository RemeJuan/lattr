package domain

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/RemeJuan/lattr/utils/error_formats"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/lib/pq"
)

var (
	TokenRepo TokenRepoInterface = &tokenRepo{}
)

var (
	queryCreateToken = "INSERT INTO tokens(Name, Token, Scopes, ExpiresAt, CreatedAt, Modified)  VALUES($1, $2, $3, $4, $5, $6) RETURNING id;"
	queryGetToken    = "SELECT name, token, scopes, expiresAt, createdAt, Modified FROM tokens WHERE id=$1;"
	queryListTokens  = "SELECT * FROM tokens"
	queryResetToken  = "UPDATE tokens SET token=$2, expiresAt=$3 modified=$4 WHERE id=$1;"
	queryDeleteToken = "DELETE FROM tokens where id=$1"
)

type TokenRepoInterface interface {
	Initialize() *sql.DB
	Create(*Token) (*Token, error_utils.MessageErr)
	Get(int64) (*Token, error_utils.MessageErr)
	List() ([]Token, error_utils.MessageErr)
	Reset(*Token) (*Token, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
}

type tokenRepo struct {
	db *sql.DB
}

func InitTokenRepository(db *sql.DB) TokenRepoInterface {
	return &tokenRepo{
		db: db,
	}
}

func (tr *tokenRepo) Initialize() *sql.DB {
	var err error
	tr.db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	checkError(err)

	fmt.Println("Connected!")

	return tr.db
}

func (tr *tokenRepo) Create(token *Token) (*Token, error_utils.MessageErr) {
	var tk int64
	stmt, err := tr.db.Prepare(queryCreateToken)

	if err != nil {
		message := fmt.Sprintf("Error when trying to prepare all entries: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}
	defer stmt.Close()

	insertResult, createErr := stmt.Query(token.Name, token.Token, pq.Array(token.Scopes), token.ExpiresAt, token.CreatedAt, token.Modified)
	if createErr != nil {
		return nil, error_formats.ParseError(createErr)
	}

	insertResult.Next()
	inErr := insertResult.Scan(&tk)

	if inErr != nil {
		message := fmt.Sprintf("error when trying to save data: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}

	token.Id = tk
	return token, nil
}

func (tr *tokenRepo) Get(id int64) (*Token, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryGetToken)

	if err != nil {
		message := fmt.Sprintf("Error retrieving record: %s", err)
		return nil, error_utils.InternalServerError(message)
	}

	defer stmt.Close()

	var token Token
	result := stmt.QueryRow(id)

	if getError := result.Scan(&token.Id, &token.Name, &token.Token, pq.Array(&token.Scopes), &token.ExpiresAt, &token.CreatedAt, &token.Modified); getError != nil {
		return nil, error_formats.ParseError(getError)
	}

	return &token, nil
}

func (tr *tokenRepo) List() ([]Token, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryListTokens)

	if err != nil {
		return nil, error_utils.InternalServerError(fmt.Sprintf("Error when trying to prepare all entries: %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, error_formats.ParseError(err)
	}
	defer rows.Close()

	results := make([]Token, 0)

	for rows.Next() {
		var token Token
		if getError := rows.Scan(&token.Id, &token.Name, &token.Token, pq.Array(&token.Scopes), &token.ExpiresAt, &token.CreatedAt, &token.Modified); getError != nil {
			message := fmt.Sprintf("Error when trying to get message: %s", getError.Error())
			return nil, error_utils.InternalServerError(message)
		}
		results = append(results, token)
	}
	if len(results) == 0 {
		return nil, error_utils.NotFoundError("no records found")
	}
	return results, nil
}

func (tr *tokenRepo) Reset(token *Token) (*Token, error_utils.MessageErr) {
	stmt, err := tr.db.Prepare(queryResetToken)

	if err != nil {
		message := fmt.Sprintf("error when trying to prepare update: %s", err.Error())
		return nil, error_utils.InternalServerError(message)
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(token.Id, token.Token, token.ExpiresAt, token.Modified)
	if updateErr != nil {
		return nil, error_formats.ParseError(updateErr)
	}
	return token, nil
}

func (tr *tokenRepo) Delete(id int64) error_utils.MessageErr {
	stmt, err := tr.db.Prepare(queryDeleteToken)
	if err != nil {
		return error_utils.InternalServerError(fmt.Sprintf("error when trying to delete record: %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		return error_utils.InternalServerError(fmt.Sprintf("error when trying to delete record %s", err.Error()))
	}
	return nil
}
