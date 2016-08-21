package retime

import (
	"testing"
	"time"
)

func TestCompileToRegexp (t *testing.T) {
	f := New("Jan _2 15:04:05")

	dt1,_ := time.ParseInLocation("Feb 12 09:34:59","Jan _2 15:04:05",time.UTC)
	dt2,_ := f.Extract("foo Feb 12 09:34:59 bar", time.UTC)

	if ! dt2.Equal(dt1) {
		t.Error("foo")
	}
}
