package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

type TimeEntryReader interface {
	Content(time.Time) ([]string, error)
}

type TimeEntryWriter interface {
	Store(TimeEntry) error
}

type FileSystemStorage struct {
	BasePath string
	Project  string
	loader   LogLoader
}

func (fs *FileSystemStorage) Path(date time.Time) string {
	return fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", fs.BasePath, fs.Project, date.Year(), date.Year(), date.Month(), date.Day())
}

func (fs *FileSystemStorage) Store(entry TimeEntry) error {
	command := entry.Tag
	date := entry.Date
	if command != "stop" {
		command = "#" + command
	}

	var traqFile string = fs.Path(date)
	var projectDir string = path.Dir(traqFile)

	_ = os.MkdirAll(projectDir, os.ModeDir|os.ModePerm)
	file, err := os.OpenFile(traqFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(Entry(date, command, ""))

	return err
}

func (fs *FileSystemStorage) Content(date time.Time) ([]string, error) {
	return fs.loader(fs.Path(date))
}
