package main

import (
	"os"
	"path"
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
	out := CaptureStdout(func() {
		WithFakeEnv(func() {
			PrintDate("example", ContentLoader, time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC))
		})
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
	out := CaptureStdout(func() {
		WithFakeEnv(func() {
			SummarizeDate("example", ContentLoader, time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC), time.Date(1986, 9, 4, 0, 0, 0, 0, time.UTC))
		})
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
		WithFakeEnv(func() {
			EvaluateDate(ContentLoader, "example", time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC))
		})
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

func TestEntry(t *testing.T) {
	expected := `Wed Sep 3 12:00:00 +0000 1986;#test;
`

	if entry := Entry(time.Date(1986, 9, 3, 12, 0, 0, 0, time.UTC), "#test", ""); entry != expected {
		t.Errorf("got wrong entry. Expected '%v' got '%v'", expected, entry)
	}
}

func TestWriteToFile(t *testing.T) {
	startDate := time.Date(2013, 1, 3, 12, 30, 0, 0, time.UTC)
	endDate := time.Date(2013, 1, 3, 13, 30, 0, 0, time.UTC)

	WithFakeEnv(func() {
		WriteToFile("example", startDate, "test")

		filePath := FilePath("example", startDate)
		out, _ := ContentLoader(filePath)
		if len(out) != 1 {
			t.Errorf("Expected different line count. Got %v\n%v", len(out), out)
		}

		if out[0] != "Thu Jan 3 12:30:00 +0000 2013;#test;" {
			t.Errorf("Expected different first line. Got %v", out[0])
		}

		WriteToFile("example", endDate, "stop")
		out, _ = ContentLoader(filePath)

		if len(out) != 2 {
			t.Errorf("Expected different line count. Got %v", len(out))
		}
		if out[1] != "Thu Jan 3 13:30:00 +0000 2013;stop;" {
			t.Errorf("Expected different stop line. Got %v", out[1])
		}

		os.RemoveAll(path.Dir(filePath))
	})
}

func TestFilePath(t *testing.T) {
	var path string = FilePath("example", time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC))

	if path != os.Getenv("TRAQ_DATA_DIR")+"/example/1986/1986-09-03" {
		t.Errorf("FilePath = %v, want %v", path, os.Getenv("TRAQ_DATA_DIR")+"/example/1986/1986-09-03")
	}
}

func TestEmptySumFile(t *testing.T) {
	content := []string{""}
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
	content := []string{
		"Mon Oct 28 21:45:33 +0100 2013;#work;",
		"Mon Oct 28 23:24:49 +0100 2013;stop;",
	}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#work"]
		if total != 5956 || !ok {
			t.Errorf("summed['#work'] = %v, want %v", total, 5956)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}

func TestNoStopSumFile(t *testing.T) {
	content := []string{
		"Mon Oct 28 20:00:00 +0100 2013;#play;",
		"Mon Oct 28 21:45:33 +0100 2013;#work;",
		"Mon Oct 28 23:24:49 +0100 2013;stop;",
	}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#play"]
		if total != 6333 || !ok {
			t.Errorf("summed['#play'] = %v, want %v", total, 6333)
		}
		total, ok = summed["#work"]
		if total != 5956 || !ok {
			t.Errorf("summed['#work'] = %v, want %v", total, 5956)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}
func TestWithStopSumFile(t *testing.T) {
	content := []string{
		"Mon Oct 28 20:00:00 +0100 2013;#play;",
		"Mon Oct 28 21:45:33 +0100 2013;stop;",
		"Mon Oct 28 21:45:33 +0100 2013;#work;",
		"Mon Oct 28 23:24:49 +0100 2013;stop;",
	}
	var summed, error = SumFile(content)

	if error == nil {
		var total, ok = summed["#play"]
		if total != 6333 || !ok {
			t.Errorf("summed['#play'] = %v, want %v", total, 6333)
		}
		total, ok = summed["#work"]
		if total != 5956 || !ok {
			t.Errorf("summed['#work'] = %v, want %v", total, 5956)
		}
	} else {
		t.Errorf("parsing error %v", error)
	}
}
