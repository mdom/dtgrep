package main

import (
	"github.com/mdom/dtgrep/fixtime"
	"testing"
	"time"
)

type FillDateTest struct {
	argument string
	result   string
}

func TestAddYear(t *testing.T) {

	now, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	tests := []FillDateTest{
		{"0000-01-02T15:04:05Z", "2006-01-02T15:04:05Z"},
		{"0000-01-03T15:04:05Z", "2005-01-03T15:04:05Z"},
	}

	for _, v := range tests {
		argument, _ := time.Parse(time.RFC3339, v.argument)
		result, _ := time.Parse(time.RFC3339, v.result)
		if fixtime.AddYear(argument, now) != result {
			t.Error("fillDate failed", argument, result)
		}
	}

}

func TestDateRange(t *testing.T) {
	from, _ := time.Parse(time.RFC3339, "2016-05-09T10:40:00Z")
	to, _ := time.Parse(time.RFC3339, "2016-05-09T11:40:00Z")

	var s, e time.Time
	var d time.Duration

	s, e = dateRange(from, to, time.Duration(0))
	if s.String() != "2016-05-09 10:40:00 +0000 UTC" || e.String() != "2016-05-09 11:40:00 +0000 UTC" {
		t.Error("specified from and to without duration")
	}

	d, _ = time.ParseDuration("20s")
	s, e = dateRange(from, time.Time{}, d)
	if s.String() != "2016-05-09 10:40:00 +0000 UTC" || e.String() != "2016-05-09 10:40:20 +0000 UTC" {
		t.Error("specified from with duration")
	}

	d, _ = time.ParseDuration("20s")
	s, e = dateRange(time.Time{}, to, d)
	if s.String() != "2016-05-09 11:39:40 +0000 UTC" || e.String() != "2016-05-09 11:40:00 +0000 UTC" {
		t.Error("specified to with duration")
	}

}
