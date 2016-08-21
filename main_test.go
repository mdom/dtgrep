package main

import (
	"testing"
	"time"
)

type FillDateTest struct {
	argument string
	result   string
}

func TestFillDate(t *testing.T) {

	now,_ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	tests := []FillDateTest{
		{"0000-01-02T15:04:05Z", "2006-01-02T15:04:05Z"},
		{"0000-01-03T15:04:05Z", "2005-01-03T15:04:05Z"},
	}

	for _, v := range tests {
		argument,_ := time.Parse(time.RFC3339,v.argument)
		result,_ := time.Parse(time.RFC3339,v.result)
		if fillDate(argument, now) != result {
			t.Error("fillDate failed", argument, result)
		}
	}

}
