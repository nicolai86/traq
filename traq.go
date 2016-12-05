/*
Package traq implements helper methods for time tracking.
*/
package traq

import (
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
	var dates = make([]time.Time, 0)
	var date = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	for {
		dates = append(dates, date)
		date = date.Add(time.Hour * 24)
		if int(date.Month()) != month {
			break
		}
	}

	return dates
}

// LogLoader defines a simple abstraction to support different loading backends
type LogLoader func(string) ([]string, error)

// ContentLoader loads file contents
func ContentLoader(filePath string) ([]string, error) {
	content, err := ioutil.ReadFile(filePath)
	lines := strings.Split(string(content), "\n")
	if lines[len(lines)-1] == "" {
		return lines[0 : len(lines)-1], err
	}
	return lines, err
}

var stopLine = regexp.MustCompile(`;stop;`)

// RunningLoader inserts a transient stop if necessary
func RunningLoader(filePath string) ([]string, error) {
	content, err := ContentLoader(filePath)

	if err == nil {
		if stopLine.MatchString(content[len(content)-1]) {
			return content, err
		}

		var line = Entry(time.Now(), "stop")
		n := len(content)
		newContent := make([]string, n+1)
		copy(newContent, content[0:])
		newContent[n] = line

		return newContent, err
	}

	return content, err
}

// SumFile evaluates the content of a traq tracking file.
// The returned map contains every tag contained in the file as well as the
// tracked duration in seconds.
func SumFile(lines []string) (map[string]int64, error) {
	var totalled = make(map[string]int64)

	var currentTag = ""
	var currentTime time.Time

	for _, line := range lines {
		if line != "" {
			var parts = strings.Split(line, ";")

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

// Handler performs actions on projects/ loaders and timeframes
type Handler func(string, LogLoader, ...time.Time)

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

func TotalDate(project string, loader LogLoader, dates ...time.Time) map[string]int64 {
	var tags = make(map[string]int64)
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
	return tags
}

// SummarizeDate prints summary informations
func SummarizeDate(project string, loader LogLoader, dates ...time.Time) {
	var tags = TotalDate(project, loader, dates...)

	date := dates[0]
	fmt.Printf("%4d-%02d-%02d\n", date.Year(), date.Month(), date.Day())
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

// Entry generates a traq entry
func Entry(date time.Time, command string) string {
	return fmt.Sprintf("%s;%s;%s\n", date.Format("Mon Jan 2 15:04:05 -0700 2006"), command, "")
}

// WriteToFile writes a given command to a traq file, converting it into a tag
// if it's no known command.
func WriteToFile(project string, date time.Time, command string) {
	if command != "stop" {
		command = "#" + command
	}

	var traqFile = FilePath(project, date)
	var projectDir = path.Dir(traqFile)

	_ = os.MkdirAll(projectDir, os.ModeDir|os.ModePerm)
	var file, error = os.OpenFile(traqFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if error == nil {
		file.WriteString(Entry(date, command))
		file.Close()
	}
}
