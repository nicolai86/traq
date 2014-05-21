package main

import (
	"strings"
	"testing"
	"time"
)

func TestDatesInMonth(t *testing.T) {
	dates := DatesInMonth(1986, 9)

	if len(dates) != 30 {
		t.Errorf("expected 30 days in Sep 1986, got %v", len(dates))
	}

	if dates[0].Weekday() != time.Monday {
		t.Errorf("Started on a Monday, got %v", dates[0].Weekday())
	}

	if dates[len(dates)-1].Weekday() != time.Tuesday {
		t.Errorf("Ended on a Tuesday, got %v", dates[len(dates)-1].Weekday())
	}
}

func TestPrintDate(t *testing.T) {
	storage := NewFixtureFileStorage()
	out := CaptureStdout(func() {
		PrintDate(storage,
			time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC),
		)
	})

	expected :=
		`Wed Sep 03 20:00:00 +0100 1986;#birth;comment
Wed Sep 03 21:45:33 +0100 1986;#chillout;
Wed Sep 03 23:24:49 +0100 1986;stop;
%%
`
	if out != expected {
		t.Errorf("unexpected PrintDate output. Expected '%v' got '%v'", expected, out)
	}
}

func TestSummarizeDate(t *testing.T) {
	storage := NewFixtureFileStorage()

	out := CaptureStdout(func() {
		SummarizeDate(storage,
			time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC),
			time.Date(1986, 9, 4, 0, 0, 0, 0, time.UTC),
		)
	})

	expectedLines := map[string]bool{
		"#birth:1.7592":    false,
		"#chillout:3.3089": false,
		"#sleeping:1.7592": false,
	}
	for _, line := range strings.Split(out, "\n") {
		expectedLines[line] = true
	}
	for key, present := range expectedLines {
		if !present {
			t.Errorf("unexpected EvaluateDate output. Expected '%v', missing from '%v'", key, out)
		}
	}
}

func TestEvaluateDate(t *testing.T) {
	out := CaptureStdout(func() {
		storage := NewFixtureFileStorage()
		EvaluateDate(storage, time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC))
	})

	expectedLines := map[string]bool{
		"1986-09-03":       false,
		"#birth:1.7592":    false,
		"#chillout:1.6544": false,
	}
	for _, line := range strings.Split(out, "\n") {
		expectedLines[line] = true
	}
	for key, present := range expectedLines {
		if !present {
			t.Errorf("unexpected EvaluateDate output. Expected '%v', missing from '%v'", key, out)
		}
	}
}

func TestEmptySumFile(t *testing.T) {
	content := []TimeEntry{}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#work"]
		if ok {
			t.Errorf("summed['#work'] = %v, should not exist", total)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}

func TestSimpleSumFile(t *testing.T) {
	content := []TimeEntry{
		TimeEntry{time.Date(2013, 10, 28, 21, 45, 33, 0, time.FixedZone("Berlin", int(time.Hour)*1)), "#work", ""},
		TimeEntry{time.Date(2013, 10, 28, 23, 15, 33, 0, time.FixedZone("Berlin", int(time.Hour)*1)), "stop", ""},
	}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#work"]
		if total != 5400 || !ok {
			t.Errorf("summed['#work'] = %v, want %v", total, 5400)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}

func TestNoStopSumFile(t *testing.T) {
	content := []TimeEntry{
		TimeEntry{time.Date(2013, 10, 28, 20, 0, 0, 0, time.UTC), "#play", ""},
		TimeEntry{time.Date(2013, 10, 28, 21, 45, 33, 0, time.UTC), "#work", ""},
		TimeEntry{time.Date(2013, 10, 28, 23, 24, 29, 0, time.UTC), "stop", ""},
	}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#play"]
		if total != 6333 || !ok {
			t.Errorf("summed['#play'] = %v, want %v", total, 6333)
		}
		total, ok = summed["#work"]
		if total != 5936 || !ok {
			t.Errorf("summed['#work'] = %v, want %v", total, 5936)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}
func TestWithStopSumFile(t *testing.T) {
	content := []TimeEntry{
		TimeEntry{time.Date(2013, 10, 28, 20, 0, 0, 0, time.UTC), "#play", ""},
		TimeEntry{time.Date(2013, 10, 28, 21, 45, 33, 0, time.UTC), "stop", ""},
		TimeEntry{time.Date(2013, 10, 28, 21, 45, 33, 0, time.UTC), "#work", ""},
		TimeEntry{time.Date(2013, 10, 28, 23, 24, 29, 0, time.UTC), "stop", ""},
	}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#play"]
		if total != 6333 || !ok {
			t.Errorf("summed['#play'] = %v, want %v", total, 6333)
		}
		total, ok = summed["#work"]
		if total != 5936 || !ok {
			t.Errorf("summed['#work'] = %v, want %v", total, 5936)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}
