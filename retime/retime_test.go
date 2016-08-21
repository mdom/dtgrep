package retime

import (
	"testing"
	"time"
)

func TestCompileToRegexp (t *testing.T) {
	f := New("Jan _2 15:04:05", time.UTC)

	dt1,_ := time.ParseInLocation("Jan _2 15:04:05", "Feb 12 09:34:59", time.UTC)
	dt2,_ := f.Extract("foo Feb 12 09:34:59 bar")

	if ! dt2.Equal(dt1) {
		t.Error("foo")
	}
}
