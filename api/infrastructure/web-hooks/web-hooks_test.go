package web_hooks

//func TestHandleWebhook(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//
//	jsonMap := make(map[string]interface{})
//	jsonMap["data"] = "testing"
//	jsonBytes, err := json.Marshal(jsonMap)
//
//	if err != nil {
//		t.Log("error should be nil", err)
//		t.Fail()
//	}
//
//	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
//
//	HandleWebhook(c)
//}
