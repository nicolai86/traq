/*
Package traq implements helper methods for time tracking.
*/
package traq

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"path"
)

var traqPath string = os.Getenv("TRAQ_DATA_DIR")

// FilePath returns the path to a traq tracking file, taking the current
// env into account.
func FilePath(project string, date time.Time) (path string) {
	return fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", traqPath, project, date.Year(), date.Year(), date.Month(), date.Day())
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

type LogLoader func (string) ([]string, error)

func ContentLoader(filePath string) ([]string, error) {
	content, err := ioutil.ReadFile(filePath)
	return strings.Split(string(content), "\n"), err
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

func Entry(date time.Time, command string) (string) {
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

	_ = os.MkdirAll(projectDir, os.ModeDir | os.ModePerm)
	var file, error = os.OpenFile(traqFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if error == nil {
		file.WriteString(Entry(date, command))
		file.Close()
	}
}
