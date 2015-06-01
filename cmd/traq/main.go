package main

import (
	"flag"
	"time"

	"github.com/nicolai86/traq"
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

	var loader = traq.ContentLoader

	if *running {
		loader = traq.RunningLoader
	}

	if *evaluate {
		if *date == "" {
			traq.EvaluateDate(loader, *project, traq.DatesInMonth(*year, *month)...)
		} else {
			traq.EvaluateDate(loader, *project, t)
		}
		return
	}

	var command = flag.Arg(0)

	var handler = traq.PrintDate
	if *summary {
		handler = traq.SummarizeDate
	}

	if command == "" {
		if *date == "" {
			handler(*project, loader, traq.DatesInMonth(*year, *month)...)
		} else {
			handler(*project, loader, t)
		}
	} else {
		traq.WriteToFile(*project, now, command)
	}
}
