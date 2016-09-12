package fixtime

import "time"

func AddDate(dt time.Time, now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
}

func AddDateHour(dt time.Time, now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
}

func AddYear(dt time.Time, now time.Time) time.Time {
	if dt.Year() == 0 {
		dt = dt.AddDate(now.Year(), 0, 0)
		if dt.After(now) {
			dt = dt.AddDate(-1, 0, 0)
		}
	}
	return dt
}
