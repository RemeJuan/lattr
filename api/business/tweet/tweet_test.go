package tweet

import (
	"testing"
)

func TestBuildTweet(t *testing.T) {
	data := make(map[string]string)
	data["message"] = "hello"
	data["time"] = "2021-07-17 12:55:50 +0200"

	expected := Tweet{
		Message:  "hello",
		UserId:   "001",
		PostTime: "2021-07-17 10:55:50 +0000 UTC",
		Status:   Pending,
	}
	t2, err := BuildTweet(data)

	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
		return
	}

	if t2 != expected {
		t.Log("Output does not match expected", t2)
		t.Fail()
	}
}
