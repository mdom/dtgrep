# NAME

dtgrep - print lines matching a date range

# SYNOPSIS

    dtgrep --from RFC3339 --to RFC3339 --format TIME_LAYOUT syslog

# DESCRIPTION

Do you even remember how often in your life you needed to find lines in
a log in a date range? And how often you build brittle regexps in grep
to match entries spanning over an hour change?

With dtgrep you don't have to. It features

* efficient binary search on normal files
* read bzip and gzip files without external dependencies
* automatically sort files
* merge lines from different files in output stream
* do as little work as necessary
* flexible syntax to declare date ranges

# EXAMPLES

But just let me show you a few examples.

The only parameter dtgrep really needs is _format_ to tell it how to
reckognize a timestamp. In this case dtgrep matches all lines from epoch to
the time dtgrep started.

    dtgrep --format "Jan _2 15:04:05" syslog

There are also some already predefined formats you can use:

    dtgrep --format apache access.log

You can specify which timerange to print:

    dtgrep --from 2006-01-02T12:00:00 --to 2006-01-02T12:15:00 syslog

If you leave one out it either defaults to epoch or the start of the program.

    dtgrep  --to 2006-01-02T12:15:00 --format rsyslog syslog

dtgrep can also read lines from stdin, but filtering those will be
slower as you can't just seek in a pipe.  It's often more efficient to
just redirect the lines from the pipe to a file first. But nothing is
stopping you to just call dtgrep directly.

    zcat syslog.gz | dtgrep --to 2006-01-02T12:15:00
    dtgrep --to 2006-01-02T12:15:00 syslog.gz

# OPTIONS

* --from RFC3339

  Print all lines from RFC3339 inclusively. Defaults to January 1,
  year 1, 00:00:00 UTC.

* --to RFC3339

  Print all lines until RFC3339 exclusively. Default to the current time.

* --format FORMAT

  FORMAT describes how a date looks. The first date found on a line is used.

  FORMAT can either a named format or any layout supported by the [time package](https://golang.org/pkg/time/#Parse).

  Additionally, dtgrep supports named formats:

  * rsyslog "Jan \_2 15:04:05"
  * apache "02/Jan/2006:15:04:05 -0700"
  * iso3339 "2006-01-02T15:04:05Z07:00"

  This parameter defaults to _rsyslog_.

* --multiline

  Print lines without timestamp between matching lines.

* --skip-dateless

  Ignore lines without timestamp.

* --location LOCATION

  If a date has no explicit timezone, interpret it as in the given
  location. LOCATION must be a valid location from the IANA Time Zone
  database, such as "America/New\_York".

  If the name is "" or "UTC, interpret dates as UTC.

  This parameter defaults to the system's local time zone.

* --help

  Shows a short help message

# ENVIRONMENT

* GO\_DATEGREP\_DEFAULT\_FORMAT

  Overwrites the default for the _--format_ parameter. The syntax is described there.

# LIMITATION

dtgrep expects the files to be sorted. If the timestamps are not
ascending, dtgrep might be exiting before the last line in its date
range is printed.

# SEE ALSO

This is by far no new idea. Just on github you can find at least 10
programs that at least the binary search part of the problem. Search
for tgrep, timegrep or dategrep.

I wrote [dategrep](http://github.com/mdom/dategrep) without knowing
about these. What is weird, as  I distinctly remember that i searched
extensively for such a program but I must have choosen very poor terms.

Since then i never lost interest in the problem. This is another try to
find the right mix of features, easy deployment and faste results.

# COPYRIGHT AND LICENSE

Copyright 2016 Mario Domgoergen <mario@domgoergen.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
