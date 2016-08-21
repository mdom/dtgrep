package retime

import (
	"bytes"
	"regexp"
	"time"
)

type Format struct {
	regexp *regexp.Regexp
	layout string
	loc  *time.Location
}

func New(layout string, loc *time.Location) Format {
	return Format{
		regexp:   compileToRegexp(layout),
		layout: layout,
		loc:  loc,
	}
}

func (f *Format) Extract(s string) (time.Time, error) {
	match := f.regexp.FindString(s)
	return time.ParseInLocation(f.layout, match, f.loc)
}

func prefixAt(s string, index int, prefix string) bool {
	return len(s) >= index+len(prefix) && s[index:index+len(prefix)] == prefix
}

func compileToRegexp(layout string) *regexp.Regexp {
	var buffer bytes.Buffer

	l := len(layout)
	for i := 0; i < l; {
		switch {
		case prefixAt(layout, i, "January"):
			buffer.WriteString(`[A-Z][a-z]{3,9}`)
			i += 7
		case prefixAt(layout, i, "Jan"):
			buffer.WriteString(`[A-Z][A-Za-z]{2}`)
			i += 3
		case prefixAt(layout, i, "Monday"):
			buffer.WriteString(`[A-Z][a-z]{6,9}`)
			i += 6
		case prefixAt(layout, i, "Mon"):
			buffer.WriteString(`[A-Z][A-Za-z]{2}`)
			i += 3
		case prefixAt(layout, i, "MST"):
			buffer.WriteString(`[A-Z]+`)
			i += 3
		// 01, 02, 03, 04, 05, 06
		case l >= i+1 && layout[i] == '0' && '1' <= layout[i+1] && layout[i+1] <= '6':
			buffer.WriteString(`\d\d`)
			i += 2
		case prefixAt(layout, i, "15"):
			buffer.WriteString(`\d\d`)
			i += 2
		case layout[i] == '1': // month
			buffer.WriteString(`\d?\d`)
			i++
		case prefixAt(layout, i, "2006"):
			buffer.WriteString(`\d{4}`)
			i += 4
		case layout[i] == '2': // day
			buffer.WriteString(`\d`)
			i++
		case prefixAt(layout, i, "15"):
			buffer.WriteString(`_\d{4}`)
			i += 5
		case prefixAt(layout, i, "_2"): //day
			buffer.WriteString(`[ \d]\d`)
			i += 2
		case layout[i] == '3': // hour12
			buffer.WriteString(`\d?d`)
			i++
		case layout[i] == '4': // minute
			buffer.WriteString(`\d\d`)
			i++
		case layout[i] == '5': // second
			buffer.WriteString(`\d\d`)
			i++
		case prefixAt(layout, i, "PM"):
			buffer.WriteString(`(PM|AM)`)
			i += 2
		case prefixAt(layout, i, "pm"):
			buffer.WriteString(`(pm|am)`)
			i += 2
		case prefixAt(layout, i, "-070000"):
			buffer.WriteString(`[+-]\d{7}`)
			i += 7
		case prefixAt(layout, i, "-07:00:00"):
			buffer.WriteString(`[+-]\d\d:\d\d:\d\d`)
			i += 9
		case prefixAt(layout, i, "-0700"):
			buffer.WriteString(`[+-]\d{4}`)
			i += 5
		case prefixAt(layout, i, "-07:00"):
			buffer.WriteString(`[+-]\d\d:\d\d`)
			i += 6
		case prefixAt(layout, i, "-07"):
			buffer.WriteString(`[+-]\d{2}`)
			i += 3
		case prefixAt(layout, i, "Z070000"):
			buffer.WriteString(`(Z|[+-]\d{7})`)
			i += 7
		case prefixAt(layout, i, "Z07:00:00"):
			buffer.WriteString(`(Z|[+-]\d\d:\d\d:\d\d)`)
			i += 9
		case prefixAt(layout, i, "Z0700"):
			buffer.WriteString(`(Z|[+-]\d{4})`)
			i += 5
		case prefixAt(layout, i, "Z07:00"):
			buffer.WriteString(`(Z|[+-]\d\d:\d\d)`)
			i += 6
		case prefixAt(layout, i, "Z07"):
			buffer.WriteString(`(Z|[+-]\d{2})`)
			i += 3
		case l >= i+1 && (layout[i:i+1] == ".0" || layout[i:i+1] == ".9"):
			j := i + 1
			ch := layout[i+2]
			for j < l && layout[j] == ch {
				j++
			}
			buffer.WriteString(`\.\d+`)
			i = i + j + 1
		default:
			buffer.WriteString(regexp.QuoteMeta(string(layout[i])))
			i++
		}
	}
	return regexp.MustCompile(buffer.String())
}
