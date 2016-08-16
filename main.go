package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"time"
	"sort"
)

var now = time.Now()
var epoch time.Time

type Format struct {
	regexp   string
	name     string
	template string
}

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

var formats = []Format{
	{
		regexp:   `^[A-Z][a-z]{2}  ?\d\d? \d{2}:\d{2}:\d{2}`,
		name:     "rsyslog",
		template: "Jan _2 15:04:05",
	},
	{
		regexp:   `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d)?(Z|[+-]\d\d:\d\d)`,
		name:     "rfc3339",
		template: time.RFC3339,
	},
}

func parse_date(date string, template string) (time.Time, error) {
	if date == "now" {
		return time.Now(), nil
	}
	dt, err := time.ParseInLocation(template, date, time.Local)
	if err != nil {
		return dt, err
	}
	now := time.Now()
	if dt.Year() == 0 {
		dt = dt.AddDate(now.Year(), 0, 0)
	}
	return dt, nil
}

func main() {

	log.SetFlags(0)
	log.SetPrefix("")

	var from_arg, to_arg, formatName string

	var options Options

	flag.StringVar(&from_arg, "from", epoch.Format(time.RFC3339), "Print all lines from `DATESPEC` inclusively.")
	flag.StringVar(&to_arg, "to", "now", "Print all lines until `DATESPEC` exclusively.")
	flag.StringVar(&formatName, "format", "rsyslog", "Use `Format` to parse file.")
	flag.BoolVar(&options.skipDateless, "skip-dateless", false, "Ignore all lines without timestamp.")
	flag.BoolVar(&options.multiline, "multiline", false, "Print all lines between the start and end line even if they are not timestamped.")

	flag.Parse()

	var err error
	options.from, err = parse_date(from_arg, time.RFC3339)
	if err != nil {
		log.Fatalln("Can't parse --from:", err)
	}
	options.to, err = parse_date(to_arg, time.RFC3339)
	if err != nil {
		log.Fatalln("Can't parse --to:", err)
	}

	var format Format
	for _, f := range formats {
		if f.name == formatName {
			format = f
			break
		}
	}

	if (format == Format{}) {
		log.Fatalln("Unknown format:", formatName)
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

func (i *Iterator) Print(to time.Time, options Options, format Format) {
	for {
		i.Line, i.Err = readline(i.Scanner)
		if i.Err == io.EOF {
			return
		}
		if i.Err != nil {
			// what file?
			log.Fatalln("Error reading file:", i.Err)
		}
		i.Time, i.Err = getLineTime(i.Line, format)

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

func getLineTime(line string, format Format) (time.Time, error) {
	re := regexp.MustCompile(format.regexp)
	time_string := re.FindString(line)
	dt, err := parse_date(time_string, format.template)
	return dt, err
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

func (i *Iterator) Scan(from, to time.Time, ignoreError bool, format Format) {
	for {
		i.Line, i.Err = readline(i.Scanner)
		if i.Err != nil {
			break
		}
		i.Time, i.Err = getLineTime(i.Line, format)
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

func findStartSeekable(f *os.File, options Options, format Format) (*bufio.Scanner, error) {

	// find block size
	block_size := int64(4096)

	file_info, err := f.Stat()
	if err != nil {
		return &bufio.Scanner{}, err
	}
	size := file_info.Size()
	min := int64(0)
	max := size / block_size
	var mid int64

	var ignore_errors = options.skipDateless || options.multiline

	for max-min > 1 {
		mid = (max + min) / 2
		f.Seek(mid*block_size, os.SEEK_SET)
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

			dt, err = getLineTime(line, format)
			if err != nil && ignore_errors {
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

	min = min * block_size
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
