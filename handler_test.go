package mockgo

import "testing"

func TestGetHandlerList(t *testing.T) {
	result := GetHandler("testdata/api", false)
	if len(result) != 1 {
		t.Error("GetHandlerList fail")
	}
}
