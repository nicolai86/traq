package main

// change traq env to use our fixtures
import (
	"bytes"
	"io"
	"os"
)

func WithFakeEnv(block func()) {
	oldEnv := os.Getenv("TRAQ_DATA_DIR")
	path, _ := os.Getwd()
	os.Setenv("TRAQ_DATA_DIR", path+"/fixtures")

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
