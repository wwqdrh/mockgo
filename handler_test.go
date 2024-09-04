package mockgo

import "testing"

func TestGetHandlerList(t *testing.T) {
	result := GetHandler("testdata/api")
	if len(result) != 1 {
		t.Error("GetHandlerList fail")
	}
}
