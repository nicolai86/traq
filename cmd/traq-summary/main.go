package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nicolai86/traq"
)

type JSONDate time.Time

func (t *JSONDate) UnmarshalJSON(data []byte) error {
	date, err := time.Parse("\"2006-01-02\"", string(data))
	if err != nil {
		return err
	}

	*t = JSONDate(date.UTC())
	return nil
}

var (
	month   = flag.Int("m", -1, "summary for a given month")
	year    = flag.Int("y", -1, "summary for a given year")
	project = flag.String("p", "timestamps", "summary for a given project")
)

func leaveFile(year int) (io.Reader, error) {
	filePath := fmt.Sprintf("%s/leave-%4d.json", os.Getenv("TRAQ_DATA_DIR"), year)
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}
	return os.Open(filePath)
}

type Leave struct {
	Holidays []JSONDate `json:"holidays"`
	Vacation []JSONDate `json:"vacation"`
	Sick     []JSONDate `json:"sick"`
}

func (l *Leave) isHoliday(date time.Time) bool {
	for _, d := range l.Holidays {
		if time.Time(d).Equal(date) {
			return true
		}
	}
	return false
}

func (l *Leave) isVacation(date time.Time) bool {
	for _, d := range l.Vacation {
		if time.Time(d).Equal(date) {
			return true
		}
	}
	return false
}

func (l *Leave) isSickDay(date time.Time) bool {
	for _, d := range l.Sick {
		if time.Time(d).Equal(date) {
			return true
		}
	}
	return false
}

func (l *Leave) expectedWorkTime(startDate, endDate time.Time) int64 {
	expected := int64(0)

	date := startDate
	for date.Before(endDate) {
		switch date.Weekday() {
		case time.Monday:
			fallthrough
		case time.Tuesday:
			fallthrough
		case time.Wednesday:
			fallthrough
		case time.Thursday:
			fallthrough
		case time.Friday:
			if !(l.isHoliday(date) || l.isVacation(date) || l.isSickDay(date)) {
				expected = expected + 8*60*60 + 60*60
			}
		}
		date = date.Add(24 * time.Hour)
	}
	return expected
}

func summarize(m map[string]int64) int64 {
	total := int64(0)
	for _, v := range m {
		total = total + v
	}
	return total
}

func main() {
	flag.Parse()

	var now = time.Now()
	if *month == -1 {
		*month = int(now.Month())
	}
	if *year == -1 {
		*year = now.Year()
	}

	r, err := leaveFile(*year)
	if err != nil {
		log.Fatal(err)
	}

	leave := &Leave{}
	json.NewDecoder(r).Decode(leave)

	startDate := time.Date(*year, time.Month(*month), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(*year, time.Month(*month+1), -1, 0, 0, 0, 0, time.UTC)

	expectedWorkTime := leave.expectedWorkTime(startDate, endDate)
	totalled := summarize(traq.TotalDate(*project, traq.ContentLoader, traq.DatesInMonth(*year, *month)...))

	fmt.Printf("work time for %4d-%02d: %2.2f vs %2.2f\n", *year, *month, float64(expectedWorkTime)/60.0/60.0, float64(totalled)/60.0/60.0)
}
