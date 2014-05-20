// Copyright 2013-2014 Raphael Randschau

/*
Package traq implements a CLI for time tracking.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type TimeEntryStorage interface {
	Store(TimeEntry) error
	Content(time.Time) ([]string, error)
}

type FileSystemStorage struct {
	BasePath string
	Project  string
	loader   LogLoader
}

func (fs *FileSystemStorage) Path(date time.Time) string {
	return fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", fs.BasePath, fs.Project, date.Year(), date.Year(), date.Month(), date.Day())
}

func (fs *FileSystemStorage) Store(entry TimeEntry) error {
	WriteToFile(fs.Project, entry.Date, entry.Tag)
	return nil
}

func (fs *FileSystemStorage) Content(date time.Time) ([]string, error) {
	return fs.loader(fs.Path(date))
}

// FilePath returns the path to a traq tracking file, taking the current
// env into account.
func FilePath(project string, date time.Time) (path string) {
	return fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", os.Getenv("TRAQ_DATA_DIR"), project, date.Year(), date.Year(), date.Month(), date.Day())
}

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
func SumFile(lines []string) (map[string]int64, error) {
	var totalled map[string]int64 = make(map[string]int64)

	var currentTag string = ""
	var currentTime time.Time

	for _, line := range lines {
		if line != "" {
			var parts []string = strings.Split(line, ";")

			var t, error = time.Parse("Mon Jan 2 15:04:05 -0700 2006", parts[0])
			if error == nil {
				if parts[1] == "" {
					currentTag = parts[1]
					totalled[currentTag] = 0
				} else if parts[1] == "stop" {
					var diff = t.Unix() - currentTime.Unix()
					totalled[currentTag] = totalled[currentTag] + diff
					currentTag = ""
				} else if currentTag != parts[1] {
					var diff = t.Unix() - currentTime.Unix()
					totalled[currentTag] = totalled[currentTag] + diff
					currentTag = parts[1]
				}

				currentTime = t
			} else {
				return totalled, error
			}
		}
	}

	delete(totalled, "")

	return totalled, nil
}

type TraqHandler func(string, LogLoader, ...time.Time)

// PrintDate prints the content of a single traqfile, identified by the project identifer
// and its date
func PrintDate(project string, loader LogLoader, dates ...time.Time) {
	for _, date := range dates {
		content, err := loader(FilePath(project, date))

		if err == nil {
			fmt.Print(strings.Join(content, "\n"))
			fmt.Println("\n%%")
		}
	}
}

func SummarizeDate(project string, loader LogLoader, dates ...time.Time) {
	var tags map[string]int64 = make(map[string]int64)
	for _, date := range dates {
		content, err := loader(FilePath(project, date))

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
func evaluateDate(contentLoader LogLoader, storage TimeEntryStorage, dates ...time.Time) {
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

func Entry(date time.Time, command string, comment string) string {
	return fmt.Sprintf("%s;%s;%s\n", date.Format("Mon Jan 2 15:04:05 -0700 2006"), command, comment)
}

// WriteToFile writes a given command to a traq file, converting it into a tag
// if it's no known command.
func WriteToFile(project string, date time.Time, command string) {
	if command != "stop" {
		command = "#" + command
	}

	var traqFile string = FilePath(project, date)
	var projectDir string = path.Dir(traqFile)

	_ = os.MkdirAll(projectDir, os.ModeDir|os.ModePerm)
	var file, error = os.OpenFile(traqFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if error == nil {
		file.WriteString(Entry(date, command, ""))
		file.Close()
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
			evaluateDate(loader, &storageProvider, DatesInMonth(*year, *month)...)
		} else {
			evaluateDate(loader, &storageProvider, t)
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
			handler(*project, loader, DatesInMonth(*year, *month)...)
		} else {
			handler(*project, loader, t)
		}
	} else {
		WriteToFile(*project, now, command)
	}
}
