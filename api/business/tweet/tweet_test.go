package tweet

import (
	"testing"

	"github.com/jinzhu/gorm"
)

func TestConvert(t *testing.T) {
	data := make(map[string]string)
	data["message"] = "hello"
	data["time"] = "2021-07-17 12:55:50 +0200"

	expected := Tweet{
		Model:    gorm.Model{},
		Message:  "hello",
		UserId:   "001",
		PostTime: "2021-07-17 10:55:50 +0000 UTC",
		Status:   Pending,
	}
	err, t2 := BuildTweet(data)

	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
		return
	}

	if t2 != expected {
		t.Log("Output does not match expected", err, t2)
		t.Fail()
	}
}
