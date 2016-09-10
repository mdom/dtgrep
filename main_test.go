package main

import (
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
		if addYear(argument, now) != result {
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

func TestDateFlagSet(t *testing.T) {
	time.Local = time.UTC
	now, _ = time.Parse(time.RFC3339, "2016-05-09T10:40:00Z")
	d := &dateFlag{}

	var err error

	err = d.Set("12")
	if err != nil || d.String() != "2016-05-09 10:12:00 +0000 UTC" {
		t.Error("Passing 12 failed")
	}

	err = d.Set("12:15")
	if err != nil || d.String() != "2016-05-09 12:15:00 +0000 UTC" {
		t.Error("Passing 12:15 failed")
	}

	err = d.Set("2016-05-08 12:15")
	if err != nil || d.String() != "2016-05-08 12:15:00 +0000 UTC" {
		t.Error("Passing 2016-05-08 12:15 failed")
	}

	err = d.Set("12:15 truncate 1h")
	if err != nil || d.String() != "2016-05-09 12:00:00 +0000 UTC" {
		t.Error("Passing 12:15 truncate 1h failed")
	}

	err = d.Set("12:15 add 5m")
	if err != nil || d.String() != "2016-05-09 12:20:00 +0000 UTC" {
		t.Error("Passing 12:15 add 5m failed")
	}

	err = d.Set("now")
	if err != nil || d.String() != "2016-05-09 10:40:00 +0000 UTC" {
		t.Error("Passing now failed")
	}

	err = d.Set("truncate 1h")
	if err != nil || d.String() != "2016-05-09 10:00:00 +0000 UTC" {
		t.Error("Passing truncate 1h failed")
	}

	err = d.Set("12:15 minus 5m")
	if err == nil || d.String() == "2016-05-09 12:10:00 +0000 UTC" {
		t.Error("Passing unknown operator minus succeeded")
	}

	err = d.Set("12:15 truncate 1h minus 5m")
	if err == nil || d.String() == "2016-05-09 11:55:00 +0000 UTC" {
		t.Error("Passing 12:15 truncate 1h minus 5m succeeded")
	}

	err = d.Set("12:15 truncate")
	if err == nil {
		t.Error("Passing 12:15 truncate without argument succeeded")
	}

	err = d.Set("12:15 truncate hour")
	if err == nil {
		t.Error("Passing 12:15 truncate hour succeeded")
	}

}
