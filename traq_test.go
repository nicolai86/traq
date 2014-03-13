package main

import (
	"os"
	"testing"
	"time"
	"bytes"
	"io"
)

// change traq env to use our fixtures
func WithFakeEnv(block func()) {
  oldEnv := os.Getenv("TRAQ_DATA_DIR")
  path, _ := os.Getwd()
  os.Setenv("TRAQ_DATA_DIR", path + "/fixtures")

  block()

  os.Setenv("TRAQ_DATA_DIR", oldEnv)
}

// capture output written to os.Stdout and return it
func CaptureStdout(block func()) string {
  old := os.Stdout // keep backup of the real stdout
  r, w, _ := os.Pipe()
  os.Stdout = w

  block()

  outC := make(chan string)
  // copy the output in a separate goroutine so printing can't block indefinitely
  go func() {
      var buf bytes.Buffer
      io.Copy(&buf, r)
      outC <- buf.String()
  }()

  // back to normal state
  w.Close()
  os.Stdout = old

  return <-outC
}

func TestDatesInMonth(t *testing.T) {
  dates := DatesInMonth(1986, 9)

  if len(dates) != 30 {
    t.Errorf("expected 30 days in Sep 1986, got %v", len(dates))
  }

  if dates[0].Weekday() != time.Monday {
    t.Errorf("Started on a Monday, got %v", dates[0].Weekday())
  }

  if dates[len(dates) - 1].Weekday() != time.Tuesday {
    t.Errorf("Ended on a Tuesday, got %v", dates[len(dates) - 1].Weekday())
  }
}

func TestPrintDate(t *testing.T) {
  out := CaptureStdout(func() {
    WithFakeEnv(func() {
      PrintDate("example", time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC))
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
