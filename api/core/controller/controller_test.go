package controller

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
)

type MockDBController struct {
	*DBController
	db *gorm.DB
}

func InitTests(t *testing.T) (*MockDBController, sqlmock.Sqlmock) {
	var controller *MockDBController
	var mock sqlmock.Sqlmock
	var err error
	var db *sql.DB

	db, mock, err = sqlmock.New()
	require.NoError(t, err)

	gbd, gErr := gorm.Open("postgres", db)
	require.NoError(t, gErr)

	controller = &MockDBController{db: gbd}

	return controller, mock
}

func TestDBController_WebHook(t *testing.T) {
	//gin.SetMode(gin.TestMode)
	//w := httptest.NewRecorder()
	//c, _ := gin.CreateTestContext(w)
	//
	//con, db := InitTests(t)

}

func TestPingDB(t *testing.T) {
	con, sqlMock := InitTests(t)

	PingDB(con.db)
	// we make sure that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestErrorCheck(t *testing.T) {
	ErrorCheck(nil)
}

func TestInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalServerError(c, errors.New(""))
}

func TestBadRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequestError(c, errors.New(""))
}
