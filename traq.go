/*
Package traq implements helper methods for time tracking.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

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

type LogLoader func(string) ([]string, error)

func ContentLoader(filePath string) ([]string, error) {
	content, err := ioutil.ReadFile(filePath)
	return strings.Split(string(content), "\n"), err
}

var stopLine = regexp.MustCompile(`;stop;`)

func RunningLoader(filePath string) ([]string, error) {
	content, err := ContentLoader(filePath)

	if err == nil {
		if stopLine.MatchString(content[len(content)-1]) {
			return content, err
		}

		var line = Entry(time.Now(), "stop")
		n := len(content)
		newContent := make([]string, n+1)
		copy(newContent, content)
		newContent[n] = line

		return newContent, err
	}

	return content, err
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

type TraqHandler func(string, ...time.Time)

// PrintDate prints the content of a single traqfile, identified by the project identifer
// and its date
func PrintDate(project string, dates ...time.Time) {
	for _, date := range dates {
		content, err := ioutil.ReadFile(FilePath(project, date))

		if err == nil {
			fmt.Print(string(content))
			fmt.Println("%%")
		}
	}
}

func SummarizeDate(project string, dates ...time.Time) {
	var tags map[string]int64 = make(map[string]int64)
	for _, date := range dates {
		content, err := ioutil.ReadFile(FilePath(project, date))

		if err == nil {
			var totalled, _ = SumFile(strings.Split(string(content), "\n"))

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

	fmt.Printf("%4d-%02d-%02d\n", year, month, day)
	for key, value := range tags {
		fmt.Printf("%s:%2.4f\n", key, float64(value)/60.0/60.0)
	}
}

// EvaluateDate prints the evaluation of a single traqfile, identified by the project identifier
// and its date
func EvaluateDate(contentLoader LogLoader, project string, dates ...time.Time) {
	for _, date := range dates {
		var content, error = contentLoader(FilePath(project, date))

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

func Entry(date time.Time, command string) string {
	return fmt.Sprintf("%s;%s;%s\n", date.Format("Mon Jan 2 15:04:05 -0700 2006"), command, "")
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
		file.WriteString(Entry(date, command))
		file.Close()
	}
}

var (
	month    int
	year     int
	day      int
	project  string = "timestamps"
	date     string
	evaluate bool
	running  bool
	summary  bool
)

func main() {
	flag.BoolVar(&evaluate, "e", false, "evaluate tracked times")
	flag.BoolVar(&running, "r", false, "add fake stop entry to evaluate if stop is missing")
	flag.BoolVar(&summary, "s", false, "summaries the given timeframe")

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

	var loader LogLoader = ContentLoader

	if running {
		loader = RunningLoader
	}

	if evaluate {
		if date == "" {
			EvaluateDate(loader, project, DatesInMonth(year, month)...)
		} else {
			EvaluateDate(loader, project, t)
		}
		return
	}

	var command string = flag.Arg(0)

	var handler TraqHandler = PrintDate
	if summary {
		handler = SummarizeDate
	}

	if command == "" {
		if date == "" {
			handler(project, DatesInMonth(year, month)...)
		} else {
			handler(project, t)
		}
	} else {
		WriteToFile(project, now, command)
	}
}
