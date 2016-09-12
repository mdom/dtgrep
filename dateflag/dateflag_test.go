package dateflag

import (
	"testing"
	"time"
)

func TestDateFlagSet(t *testing.T) {
	time.Local = time.UTC
	now, _ := time.Parse(time.RFC3339, "2016-05-09T10:40:00Z")
	d := &DateFlag{Now: now}

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
