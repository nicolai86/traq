// Copyright 2013-2014 Raphael Randschau

/*
Package traq implements a CLI for time tracking.
*/
package main

import (
	"flag"
	"os"
	"time"
)

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

	printDate := !*summary
	handler := func(storage FileSystemStorage, dates ...time.Time) {
		if printDate {
			PrintDate(&storageProvider, dates...)
		} else {
			PrintSummary(SummarizeDate(&storageProvider, dates...))
		}
	}

	if command == "" {
		if *date == "" {
			handler(storageProvider, DatesInMonth(*year, *month)...)
		} else {
			handler(storageProvider, t)
		}
	} else {
		storageProvider.Store(TimeEntry{now, command, ""})
	}
}
