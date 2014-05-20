package main

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestRunningLoader(t *testing.T) {
	startDate := time.Date(2013, 1, 3, 12, 30, 0, 0, time.UTC)

	tags := map[string]bool{
		"#test": false,
		"stop":  false,
	}

	WithFakeEnv(func() {
		WriteToFile("example", startDate, "test")

		filePath := FilePath("example", startDate)
		out, _ := RunningLoader(filePath)

		for _, line := range out {
			entry := ReadEntry(line)
			tags[entry.Tag] = true
		}

		if len(out) != 3 {
			t.Errorf("Expected different line count. Got %v", len(out))
		}

		for key, present := range tags {
			if !present {
				t.Errorf("unexpected EvaluateDate output. Expected '%v', missing from '%v'", key, out)
			}
		}

		os.RemoveAll(path.Dir(filePath))
	})
}
