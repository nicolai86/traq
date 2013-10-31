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
