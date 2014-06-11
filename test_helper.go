package traq

// change traq env to use our fixtures
import (
	"bytes"
	"io"
	"os"
)

func NewFixtureFileStorage() *FileSystemStorage {
	path, err := os.Getwd()
	if err != nil {
		return nil
	}
	return &FileSystemStorage{path + "/fixtures", "example", ContentLoader}
}
func NewRunningFixtureFileStorage() *FileSystemStorage {
	path, err := os.Getwd()
	if err != nil {
		return nil
	}
	return &FileSystemStorage{path + "/fixtures", "example", RunningLoader}
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
