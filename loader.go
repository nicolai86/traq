package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

type TimeEntry struct {
	Date    time.Time
	Tag     string
	Comment string
}

func ReadEntry(line string) TimeEntry {
	entry := TimeEntry{}

	parts := strings.Split(line, ";")
	if len(parts) != 3 {
		return entry
	}

	var err error
	entry.Date, err = time.Parse("Mon Jan 2 15:04:05 -0700 2006", parts[0])
	if err != nil {
		return entry
	}

	entry.Tag = parts[1]
	entry.Comment = parts[2]
	return entry
}

type LogLoader func(string) ([]string, error)

func ContentLoader(filePath string) ([]string, error) {
	content, err := ioutil.ReadFile(filePath)
	lines := strings.Split(string(content), "\n")
	if lines[len(lines)-1] == "" {
		return lines[0 : len(lines)-1], err
	}
	return lines, err
}

var stopLine = regexp.MustCompile(`;stop;`)

func RunningLoader(filePath string) ([]string, error) {
	content, err := ContentLoader(filePath)

	if err != nil {
		return content, err
	}

	if stopLine.MatchString(content[len(content)-1]) {
		return content, err
	}

	var line = fmt.Sprintf("%s;stop;\n", time.Now().Format("Mon Jan 2 15:04:05 -0700 2006"))

	n := len(content)
	newContent := make([]string, n+1)
	copy(newContent, content)
	newContent[n] = line

	return newContent, err
}
