package format

import (
	"regexp"
	"bytes"
)

type Format struct {
	regexp   string
	name     string
	template string
}

func CompileFormatToRegexp(layout string) string {
	var buffer bytes.Buffer

	l := len(layout)
	for i := 0; i < l; {
		switch {
		case l >= i+7 && layout[i:i+7] == "January":
			buffer.WriteString(`[A-Z][a-z]{3,9}`)
			i += 7
		case l >= i+3 && layout[i:i+3] == "Jan":
			buffer.WriteString(`[A-Z][A-Za-z]{2}`)
			i += 3
		case l >= i+6 && layout[i:i+6] == "Monday":
			buffer.WriteString(`[A-Z][a-z]{6,9}`)
			i += 6
		case l >= i+3 && layout[i:i+3] == "Mon":
			buffer.WriteString(`[A-Z][A-Za-z]{2}`)
			i += 3
		case l >= i+3 && layout[i:i+3] == "MST":
			buffer.WriteString(`[A-Z]+`)
			i += 3
		// 01, 02, 03, 04, 05, 06
		case l >= i+1 && layout[i] == '0' && '1' <= layout[i+1] && layout[i+1] <= '6':
			buffer.WriteString(layout[i:i+1])
			i += 2
		case l >= i+1 && layout[i:i+1] == "15": //hour
			buffer.WriteString(`\d\d`)
			i += 2
		case layout[i] == '1': // month
			buffer.WriteString(`\d?\d`)
			i++
		case l >= i+4 && layout[i:i+4] == "2006": //year
			buffer.WriteString(`\d{4}`)
			i += 4
		case layout[i] == '2': // day
			buffer.WriteString(`\d`)
			i++
		case l >= i+5 && layout[i:i+5] == "_2006":
			buffer.WriteString(`_\d{4}`)
			i += 5
		case l >= i+2 && layout[i:i+1] == "_2": // day
			buffer.WriteString(`[ \d]\d`)
			i += 2
		case  layout[i] == '3': // hour12
			buffer.WriteString(`\d?\d`)
			i++
		case layout[i] == '4': // minute
			buffer.WriteString(`\d?\d`)
			i++
		case layout[i] == '5': // second
			buffer.WriteString(`\d?\d`)
			i++
		case l > i+1 && layout[i:i+1] == "PM":
			buffer.WriteString(`(PM|AM)`)
			i += 2
		case l > i+1 && layout[i:i+1] == "pm":
			buffer.WriteString(`(pm|am)`)
			i += 2
		case l >= i+7 && layout[i:i+7] == "-070000":
			buffer.WriteString(`[+-]\d{7}`)
			i += 7
		case l >= i+9 && layout[i:i+9] == "-07:00:00":
			buffer.WriteString(`[+-]\d\d:\d\d:\d\d`)
			i += 9
		case l >= i+5 && layout[i:i+5] == "-0700":
			buffer.WriteString(`[+-]\d{4}`)
			i += 5
		case l >= i+6 && layout[i:i+6] == "-07:00" :
			buffer.WriteString(`[+-]\d\d:\d\d`)
			i += 6
		case l >= i+3 && layout[i:i+3] == "-07" :
			buffer.WriteString(`[+-]\d{2}`)
			i += 3
		case l >= i+7 && layout[i:i+7] == "Z070000":
			buffer.WriteString(`(Z|[+-]\d{7})`)
			i += 7
		case l >= i+9 && layout[i:i+9] == "Z07:00:00":
			buffer.WriteString(`(Z|[+-]\d\d:\d\d:\d\d)`)
			i += 9
		case l >= i+5 && layout[i:i+5] == "Z0700":
			buffer.WriteString(`(Z|[+-]\d{4})`)
			i += 5
		case l >= i+6 && layout[i:i+6] == "Z07:00":
			buffer.WriteString(`(Z|[+-]\d\d:\d\d)`)
			i += 6
		case l >= i+3 && layout[i:i+3] == "Z07":
			buffer.WriteString(`(Z|[+-]\d{2})`)
			i += 3
		case l >= i+1 && ( layout[i:i+2] == ".0" || layout[i:i+2] == ".9"):
			j := i + 1
			ch := layout[i+2]
			for j < l && layout[j] == ch {
				j++
			}
			buffer.WriteString(`\.\d+`)
			i = i + j +1
		default:
			buffer.WriteString(regexp.QuoteMeta(string(layout[i])))
			i++
		}
	}
	return buffer.String()
}
