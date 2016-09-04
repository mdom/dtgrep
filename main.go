package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"github.com/mdom/dtgrep/retime"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

var now = time.Now()
var epoch time.Time
var loc = time.Local

var Version = "unknown"
var CommitHash = "unknown"
var BuildDate = "unknown"

type Options struct {
	from, to     time.Time
	skipDateless bool
	multiline    bool
}

type Iterator struct {
	filename string
	reader   io.Reader
	*bufio.Scanner
	Line string
	Time time.Time
	Err  error
}

type Formats struct {
	template string
	complete func(date time.Time, now time.Time) time.Time
}

type Iterators []*Iterator

func (it Iterators) Len() int           { return len(it) }
func (it Iterators) Swap(i, j int)      { it[i], it[j] = it[j], it[i] }
func (it Iterators) Less(i, j int) bool { return it[i].Time.Before(it[j].Time) }

func inTimeRange(s *Iterator, from, to time.Time) bool {
	dt := s.Time
	return (dt.Equal(from) || dt.After(from)) && dt.Before(to)
}

func filter(s Iterators, from, to time.Time) Iterators {
	var p Iterators
	for _, v := range s {
		if v.Err == nil && inTimeRange(v, from, to) {
			p = append(p, v)
		}
	}
	return p
}

var formats = map[string]string{
	"rsyslog": "Jan _2 15:04:05",
	"rfc3339": time.RFC3339,
	"apache":  "02/Jan/2006:15:04:05 -0700",
}

func parseDate(dateSpec string, template string) (time.Time, error) {

	datePart := dateSpec
	modPart := ""

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
			return time.Time{}, errors.New("Missing argument for " + fields[i])
		}

		d, err := time.ParseDuration(fields[i+1])
		if err != nil {
			return time.Time{}, err
		}
		switch fields[i] {
		case "truncate":
			modifiers = append(modifiers, func(t time.Time) time.Time { return t.Truncate(d) })
		case "add":
			modifiers = append(modifiers, func(t time.Time) time.Time { return t.Add(d) })
		default:
			return time.Time{}, errors.New("Unknown operator " + fields[i])
		}
	}

	var dt time.Time

	if datePart == "now" || datePart == "" {
		dt = time.Now()
	} else {
		specs := []Formats{
			{"15:04", addDate},
			{time.RFC3339, returnDate},
		}
		var err error
		for _, spec := range specs {
			dt, err = time.ParseInLocation(spec.template, datePart, time.Local)
			if err == nil {
				dt = spec.complete(dt, time.Now())
				break
			}
		}
		if err != nil {
			return time.Time{}, err
		}
	}
	for _, mod := range modifiers {
		dt = mod(dt)
	}
	log.Fatalln(dt)
	return dt, nil
}

func returnDate(dt time.Time, now time.Time) time.Time {
	return dt
}

func addDate(dt time.Time, now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location())
}

func addYear(dt time.Time, now time.Time) time.Time {
	if dt.Year() == 0 {
		dt = dt.AddDate(now.Year(), 0, 0)
		if dt.After(now) {
			dt = dt.AddDate(-1, 0, 0)
		}
	}
	return dt
}

func main() {

	log.SetFlags(0)
	log.SetPrefix("")

	var fromArg, toArg, formatName, location string

	var options Options

	defaultFormat := "rsyslog"
	if os.Getenv("GO_DATEGREP_FORMAT") != "" {
		defaultFormat = os.Getenv("GO_DATEGREP_FORMAT")
	}

	flag.StringVar(&fromArg, "from", epoch.Format(time.RFC3339), "Print all lines from `DATESPEC` inclusively.")
	flag.StringVar(&toArg, "to", "now", "Print all lines until `DATESPEC` exclusively.")
	flag.StringVar(&formatName, "format", defaultFormat, "Use `Format` to parse file.")
	flag.BoolVar(&options.skipDateless, "skip-dateless", false, "Ignore all lines without timestamp.")
	flag.BoolVar(&options.multiline, "multiline", false, "Print all lines between the start and end line even if they are not timestamped.")
	flag.StringVar(&location, "location", time.Local.String(), "Use location in the absence of any timezone information.")

	var displayVersion bool
	flag.BoolVar(&displayVersion, "version", false, "Display version")

	flag.Parse()

	if displayVersion {
		log.Printf("version: %s\ncommit: %s\nbuild date: %s\n",
			Version, CommitHash, BuildDate)
		return
	}

	var err error

	loc, err = time.LoadLocation(location)
	if err != nil {
		log.Fatalln("Can't load location:", err)
	}

	options.from, err = parseDate(fromArg, time.RFC3339)
	if err != nil {
		log.Fatalln("Can't parse --from:", err)
	}
	options.to, err = parseDate(toArg, time.RFC3339)
	if err != nil {
		log.Fatalln("Can't parse --to:", err)
	}

	var format retime.Format
	for name, template := range formats {
		if name == formatName {
			format, err = retime.New(template, loc)
			if err != nil {
				log.Fatalln("Can't create format:", err)
			}
			break
		}
	}

	if (format == retime.Format{}) {
		format, err = retime.New(formatName, loc)
		if err != nil {
			log.Fatalln("Can't create format:", err)
		}
	}

	var iterators = make(Iterators, 0)

	if len(flag.Args()) > 0 {
		for _, filename := range flag.Args() {

			if filename == "-" {
				i := &Iterator{filename: filename, reader: os.Stdin, Scanner: bufio.NewScanner(os.Stdin)}
				iterators = append(iterators, i)
				continue
			}

			file, err := os.Open(filename)
			if err != nil {
				log.Fatalln("Cannot open", filename, ":", err)
			}
			defer file.Close()

			// mimeType support?
			ext := path.Ext(filename)
			if ext == ".gz" || ext == ".z" {
				r, err := gzip.NewReader(file)
				defer r.Close()
				if err != nil {
					log.Fatalln("Cannot open", filename, ":", err)
				}
				i := &Iterator{filename: filename, reader: r, Scanner: bufio.NewScanner(r)}
				iterators = append(iterators, i)
			} else if ext == ".bz2" || ext == ".bz" {
				r := bzip2.NewReader(file)
				i := &Iterator{filename: filename, reader: r, Scanner: bufio.NewScanner(r)}
				iterators = append(iterators, i)
			} else {
				scanner, err := findStartSeekable(file, options, format)
				switch {
				case err == io.EOF:
					// daterange not in file, skip
					continue
				case err != nil:
					log.Fatalln("Error finding dates in ", filename, ":", err)
				}
				i := &Iterator{filename: filename, reader: file, Scanner: scanner}
				iterators = append(iterators, i)
			}
		}
	} else {
		i := &Iterator{filename: "-", reader: os.Stdin, Scanner: bufio.NewScanner(os.Stdin)}
		iterators = append(iterators, i)
	}

	var ignoreError = options.skipDateless || options.multiline
	for _, i := range iterators {
		i.Scan(options.from, options.to, ignoreError, format)
	}

	for {

		iterators = filter(iterators, options.from, options.to)
		sort.Sort(iterators)

		if len(iterators) > 0 {
			var until time.Time
			if len(iterators) > 1 {
				until = iterators[1].Time
			} else {
				until = options.to
			}
			i := iterators[0]
			fmt.Println(i.Line)
			i.Print(until, options, format)
		} else {
			break
		}
	}
}

func (i *Iterator) Print(to time.Time, options Options, format retime.Format) {
	for {
		i.Line, i.Err = readline(i.Scanner)
		if i.Err == io.EOF {
			return
		}
		if i.Err != nil {
			// what file?
			log.Fatalln("Error reading file:", i.Err)
		}
		i.Time, i.Err = format.Extract(i.Line)
		i.Time = addYear(i.Time, time.Now())

		switch {
		case i.Err != nil && options.multiline:
			fmt.Println(i.Line)
		case i.Err != nil && options.skipDateless:
			continue
		case i.Err != nil:
			log.Fatalln("Aborting. Found line without date:", i.Line)
		case i.Time.Before(to):
			fmt.Println(i.Line)
		default:
			return
		}
	}
}

func readline(s *bufio.Scanner) (string, error) {
	ret := s.Scan()
	if !ret && s.Err() == nil {
		return "", io.EOF
	}
	if !ret {
		return "", s.Err()
	}
	return s.Text(), nil
}

func (i *Iterator) Scan(from, to time.Time, ignoreError bool, format retime.Format) {
	for {
		i.Line, i.Err = readline(i.Scanner)
		if i.Err != nil {
			break
		}
		i.Time, i.Err = format.Extract(i.Line)
		i.Time = addYear(i.Time, time.Now())
		if i.Err != nil && ignoreError {
			continue
		}
		if i.Err != nil {
			log.Fatalln("Aborting. Found line without date:", i.Line)
		}
		if i.Time.After(to) {
			i.Err = io.EOF
			break
		}
		if i.Time.Equal(from) || i.Time.After(from) {
			break
		}
	}
}

func findStartSeekable(f *os.File, options Options, format retime.Format) (*bufio.Scanner, error) {

	// find block size
	blockSize := int64(4096)

	fileInfo, err := f.Stat()
	if err != nil {
		return &bufio.Scanner{}, err
	}
	size := fileInfo.Size()
	min := int64(0)
	max := size / blockSize
	var mid int64

	var ignoreErrors = options.skipDateless || options.multiline

	for max-min > 1 {
		mid = (max + min) / 2
		f.Seek(mid*blockSize, os.SEEK_SET)
		scanner := bufio.NewScanner(f)

		_, err := readline(scanner) // skip partial line
		if err != nil {
			return scanner, err
		}

		var dt time.Time

		for {
			line, err := readline(scanner)
			if err != nil {
				return scanner, err
			}

			dt, err = format.Extract(line)
			dt = addYear(dt, time.Now())
			if err != nil && ignoreErrors {
				continue
			}
			if err != nil {
				log.Fatalln("Aborting. Found line without date:", line)
			}
			break
			// optimization: while searching next line we entered next block
		}

		if dt.Before(options.from) {
			min = mid
		} else {
			max = mid
		}
	}

	min = min * blockSize
	_, err = f.Seek(min, os.SEEK_SET)
	if err != nil {
		return &bufio.Scanner{}, err
	}
	scanner := bufio.NewScanner(f)
	if min > 0 {
		_, err := readline(scanner) // skip partial line
		if err != nil {
			return scanner, err
		}
	}

	return scanner, nil
}
