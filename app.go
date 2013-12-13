package main

import (
	"flag"
	"time"
	"traq"
	"regexp"
)

var month int
var year int
var day int
var project string = "timestamps"
var date string
var evaluate bool
var running bool

var stopLine = regexp.MustCompile(`;stop;`)

func RunningLoader (filePath string) ([]string, error) {
	content, err := traq.ContentLoader(filePath)

	if err == nil {
		if stopLine.MatchString(content[len(content)-1]) {
			return content, err
		}

		var line = traq.Entry(time.Now(), "stop")
		n := len(content)
		newContent := make([]string, n + 1)
		copy(newContent, content)
		newContent[n] = line

		return newContent, err
	}

	return content, err
}

func main() {
	flag.BoolVar(&evaluate, "e", false, "evaluate tracked times")
	flag.BoolVar(&running, "r", false, "add fake stop entry to evaluate if stop is missing")

	flag.IntVar(&year, "y", 0, "print tracked times for a given year")
	flag.IntVar(&month, "m", 0, "print tracked times for a given month")

	flag.StringVar(&date, "d", "", "print tracked times for a given date")
	flag.StringVar(&project, "p", "", "print data for a given project")

	flag.Parse()

	var now = time.Now()
	var t, error = time.Parse("2006-01-02", date)
	if error == nil {
		year = t.Year()
		month = int(t.Month())
		day = t.Day()
	} else {
		if month == 0 && year == 0 {
			day = now.Day()
		} else {
			day = 1
		}
		if year == 0 {
			year = now.Year()
		}
		if month == 0 {
			month = int(now.Month())
		}
	}

	var loader traq.LogLoader = traq.ContentLoader

	if running {
		loader = RunningLoader
	}

	if evaluate {
		if date == "" {
			traq.EvaluateDate(loader, project, traq.DatesInMonth(year, month)...)
		} else {
			traq.EvaluateDate(loader, project, t)
		}
		return
	}

	var command string = flag.Arg(0)

	if command == "" {
		if date == "" {
			traq.PrintDate(project, traq.DatesInMonth(year, month)...)
		} else {
			traq.PrintDate(project, t)
		}
	} else {
		traq.WriteToFile(project, now, command)
	}
}
