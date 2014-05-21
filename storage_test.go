package main

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestFileStorageStore(t *testing.T) {
	startDate := time.Date(2013, 1, 3, 12, 30, 0, 0, time.UTC)
	endDate := time.Date(2013, 1, 3, 13, 30, 0, 0, time.UTC)

	storage := NewFixtureFileStorage()
	storage.Store(TimeEntry{startDate, "test", ""})

	out, _ := ContentLoader(storage.Path(startDate))
	if len(out) != 1 {
		t.Errorf("Expected different line count. Got %v\n%v", len(out), out)
	}

	if out[0] != "Thu Jan 3 12:30:00 +0000 2013;#test;" {
		t.Errorf("Expected different first line. Got %v", out[0])
	}

	storage.Store(TimeEntry{endDate, "stop", ""})
	out, _ = ContentLoader(storage.Path(startDate))

	if len(out) != 2 {
		t.Errorf("Expected different line count. Got %v", len(out))
	}
	if out[1] != "Thu Jan 3 13:30:00 +0000 2013;stop;" {
		t.Errorf("Expected different stop line. Got %v", out[1])
	}

	os.RemoveAll(path.Dir(storage.Path(startDate)))
}
