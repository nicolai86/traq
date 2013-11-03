package traq

import (
	"os"
	"testing"
	"time"
	"traq"
)

func TestFilePath(t *testing.T) {
	var path string = traq.FilePath("example", time.Date(1986, 9, 3, 0, 0, 0, 0, time.UTC))

	if path != os.Getenv("TRAQ_DATA_DIR")+"/example/1986/1986-09-03" {
		t.Errorf("FilePath = %v, want %v", path, os.Getenv("TRAQ_DATA_DIR")+"/example/1986/1986-09-03")
	}
}

func TestEmptySumFile(t *testing.T) {
	var content string = ""
	var summed, error = traq.SumFile(content)

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
	var content string = "Mon Oct 28 21:45:33 +0100 2013;#work;\nMon Oct 28 23:24:49 +0100 2013;stop;\n"
	var summed, error = traq.SumFile(content)

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
	var content string = "Mon Oct 28 20:00:00 +0100 2013;#play;\nMon Oct 28 21:45:33 +0100 2013;#work;\nMon Oct 28 23:24:49 +0100 2013;stop;\n"
	var summed, error = traq.SumFile(content)

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
	var content string = "Mon Oct 28 20:00:00 +0100 2013;#play;\nMon Oct 28 21:45:33 +0100 2013;stop;\nMon Oct 28 21:45:33 +0100 2013;#work;\nMon Oct 28 23:24:49 +0100 2013;stop;\n"
	var summed, error = traq.SumFile(content)

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
