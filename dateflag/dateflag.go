package dateflag

import (
	"errors"
	"github.com/mdom/dtgrep/fixtime"
	"regexp"
	"strings"
	"time"
)

type formats struct {
	template string
	complete func(date time.Time, now time.Time) time.Time
}

type DateFlag struct {
	date, Now time.Time
}

func (d *DateFlag) String() string {
	return d.date.String()
}

func (d *DateFlag) Get() time.Time {
	return d.date
}

func (d *DateFlag) Set(dateSpec string) error {

	datePart := dateSpec
	modPart := ""

	if d.Now.IsZero() {
		d.Now = time.Now()
	}

	results := regexp.MustCompile(`add|truncate`).FindStringIndex(dateSpec)
	if results != nil {
		idx := results[0]
		if idx == 0 {
			datePart = ""
			modPart = dateSpec
		} else {
			datePart = dateSpec[:idx-1]
			modPart = dateSpec[idx:]
		}
	}

	var modifiers []func(time.Time) time.Time
	fields := strings.Fields(modPart)

	for i := 0; i < len(fields); i += 2 {

		if i+1 == len(fields) {
			return errors.New("Missing argument for " + fields[i])
		}

		d, err := time.ParseDuration(fields[i+1])
		if err != nil {
			return err
		}
		switch fields[i] {
		case "truncate":
			modifiers = append(modifiers, func(t time.Time) time.Time { return t.Truncate(d) })
		case "add":
			modifiers = append(modifiers, func(t time.Time) time.Time { return t.Add(d) })
		default:
			return errors.New("Unknown operator " + fields[i])
		}
	}

	var dt time.Time

	if datePart == "now" || datePart == "" {
		dt = d.Now
	} else {
		specs := []formats{
			{"04", fixtime.AddDateHour},
			{"15:04", fixtime.AddDate},
			{"15:04:05", fixtime.AddDate},
			{"2006-01-02 15:04", nil},
			{"2006-01-02 15:04:05", nil},
			{"2006-01-02 15:04:05Z07:00", nil},
			{time.RFC3339, nil},
		}
		var err error
		for _, spec := range specs {
			dt, err = time.ParseInLocation(spec.template, datePart, time.Local)
			if err == nil {
				if spec.complete != nil {
					dt = spec.complete(dt, d.Now)
				}
				break
			}
		}
		if err != nil {
			return err
		}
	}
	for _, mod := range modifiers {
		dt = mod(dt)
	}
	d.date = dt
	return nil
}
