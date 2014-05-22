// Copyright 2013-2014 Raphael Randschau

/*
Package traq implements a CLI for time tracking.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// DatesInMonth calculates the days of a given month and year.
func DatesInMonth(year int, month int) []time.Time {
	var dates []time.Time = make([]time.Time, 0)
	var date time.Time = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	for {
		dates = append(dates, date)
		date = date.Add(time.Hour * 24)
		if int(date.Month()) != month {
			break
		}
	}

	return dates
}

// SumFile evaluates the content of a traq tracking file.
// The returned map contains every tag contained in the file as well as the
// tracked duration in seconds.
func SumFile(lines []TimeEntry) (map[string]int64, error) {
	var totalled map[string]int64 = make(map[string]int64)

	var currentTag string = ""
	var currentTime time.Time

	for _, entry := range lines {
		_, ok := totalled[entry.Tag]
		if !ok {
			totalled[entry.Tag] = 0
		}

		if entry.Tag == "stop" {
			var diff = entry.Date.Unix() - currentTime.Unix()
			totalled[currentTag] = totalled[currentTag] + diff
			currentTag = ""
		} else if currentTag != entry.Tag {
			var diff = entry.Date.Unix() - currentTime.Unix()
			totalled[currentTag] = totalled[currentTag] + diff
			currentTag = entry.Tag
		}

		currentTime = entry.Date
	}

	// we did not end with a stop tag
	// this problem can be avoided by using the RunningEvaluator when parsing a time file
	if currentTag != "stop" {
		var diff = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).Unix() - currentTime.Unix()
		totalled[currentTag] = totalled[currentTag] + diff
		currentTag = ""
	}

	delete(totalled, "")
	delete(totalled, "stop")

	return totalled, nil
}

type TraqHandler func(TimeEntryReader, ...time.Time)

// PrintDate prints the content of a single traqfile, identified by the project identifer
// and its date
func PrintDate(storage TimeEntryReader, dates ...time.Time) {
	for _, date := range dates {
		content, err := storage.Content(date)

		if err == nil {
			for _, entry := range content {
				fmt.Printf("%s;%s;%s\n", entry.Date.Format("Mon Jan 02 15:04:05 -0700 2006"), entry.Tag, entry.Comment)
			}
			fmt.Println("%%")
		}
	}
}

func SummarizeDate(storage TimeEntryReader, dates ...time.Time) {
	var tags map[string]int64 = make(map[string]int64)
	for _, date := range dates {
		content, err := storage.Content(date)

		if err == nil {
			var totalled, _ = SumFile(content)

			for key, value := range totalled {
				current, ok := tags[key]
				if !ok {
					tags[key] = value
				} else {
					tags[key] = value + current
				}
			}
		}
	}

	date := dates[0]
	fmt.Printf("%4d-%02d-%02d\n", date.Year(), date.Month(), date.Day())
	for key, value := range tags {
		fmt.Printf("%s:%2.4f\n", key, float64(value)/60.0/60.0)
	}
}

// EvaluateDate prints the evaluation of a single traqfile, identified by the project identifier
// and its date
func EvaluateDate(storage TimeEntryReader, dates ...time.Time) {
	for _, date := range dates {
		var content, error = storage.Content(date)

		if error == nil {
			fmt.Printf("%d-%02d-%02d\n", date.Year(), date.Month(), date.Day())
			var totalled, _ = SumFile(content)
			// TODO handle errors
			for key, value := range totalled {
				fmt.Printf("%s:%2.4f\n", key, float64(value)/60.0/60.0)
			}

			fmt.Println("%%")
		}
	}
}

func main() {
	var (
		month    = flag.Int("m", 0, "print tracked times for a given month")
		year     = flag.Int("y", 0, "print tracked times for a given year")
		project  = flag.String("p", "timestamps", "print data for a given project")
		date     = flag.String("d", "", "print tracked times for a given date")
		evaluate = flag.Bool("e", false, "evaluate tracked times")
		running  = flag.Bool("r", false, "add fake stop entry to evaluate if stop is missing")
		summary  = flag.Bool("s", false, "summaries the given timeframe")
	)

	flag.Parse()

	var now = time.Now()
	var t, error = time.Parse("2006-01-02", *date)
	if error == nil {
		*year = t.Year()
		*month = int(t.Month())
	} else {
		if *year == 0 {
			*year = now.Year()
		}
		if *month == 0 {
			*month = int(now.Month())
		}
	}

	var loader LogLoader = ContentLoader

	if *running {
		loader = RunningLoader
	}

	storageProvider := FileSystemStorage{os.Getenv("TRAQ_DATA_DIR"), *project, loader}

	if *evaluate {
		if *date == "" {
			EvaluateDate(&storageProvider, DatesInMonth(*year, *month)...)
		} else {
			EvaluateDate(&storageProvider, t)
		}
		return
	}

	var command string = flag.Arg(0)

	var handler TraqHandler = PrintDate
	if *summary {
		handler = SummarizeDate
	}

	if command == "" {
		if *date == "" {
			handler(&storageProvider, DatesInMonth(*year, *month)...)
		} else {
			handler(&storageProvider, t)
		}
	} else {
		storageProvider.Store(TimeEntry{now, command, ""})
	}
}
