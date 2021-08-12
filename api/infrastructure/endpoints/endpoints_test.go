package endpoints

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	const testKey = "test_key"
	const testValue = "test_value"

	c.Params = []gin.Param{{Key: testKey, Value: testValue}}

	param := GetParam(c, testKey)

	if param != testValue {
		t.Log("Function did not return expected value", param)
		t.Fail()
	}
}
