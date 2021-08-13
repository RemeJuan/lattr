package controller

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDBController struct {
	*DBController
	db *gorm.DB
}

func TestDBController_WebHook_Success(t *testing.T) {
	jsonBody := `{"data": "the title"}`
	rr := webHook(t, jsonBody)

	assert.EqualValues(t, http.StatusCreated, rr.Code)
}

func TestDBController_WebHook_BadRequest(t *testing.T) {
	jsonBody := `{"money": "the title"}`
	rr := webHook(t, jsonBody)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
}

func TestDBController_WebHook_UnprocessableEntity(t *testing.T) {
	jsonBody := ``
	rr := webHook(t, jsonBody)

	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestPingDB(t *testing.T) {
	con, sqlMock := initTests(t)

	PingDB(con.db)
	// we make sure that all expectations were met
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func initTests(t *testing.T) (*MockDBController, sqlmock.Sqlmock) {
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

func webHook(t *testing.T, jsonBody string) *httptest.ResponseRecorder {
	con, _ := initTests(t)
	gin.SetMode(gin.TestMode)

	reqBody := bytes.NewBufferString(jsonBody)

	r := gin.Default()
	req, err := http.NewRequest(http.MethodPost, "/create", reqBody)

	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	r.POST("/create", con.WebHook)
	r.ServeHTTP(rr, req)

	return rr
}