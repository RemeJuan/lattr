package controller

import (
	"database/sql"
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

func testInit(t *testing.T) {
	var controller *MockDBController
	var mock sqlmock.Sqlmock
	var err error
	var db *sql.DB

	db, mock, err = sqlmock.New()
	require.NoError(t, err)

	gbd, gErr := gorm.Open("postgress", db)
	require.NoError(t, gErr)

	controller = &MockDBController{db: gbd}
}

func (con *MockDBController) TestDBController_WebHook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testInit(t)

	con.WebHook(c)
}
