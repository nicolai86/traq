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

func serialize(entry TimeEntry) string {
	return fmt.Sprintf("%s;%s;%s\n", entry.Date.Format("Mon Jan 2 15:04:05 -0700 2006"), entry.Tag, entry.Comment)
}

func (fs *FileSystemStorage) Store(entry TimeEntry) error {
	command := entry.Tag
	date := entry.Date
	if command != "stop" {
		command = "#" + command
	}
	entry.Tag = command

	var traqFile string = fs.Path(date)
	var projectDir string = path.Dir(traqFile)

	_ = os.MkdirAll(projectDir, os.ModeDir|os.ModePerm)
	file, err := os.OpenFile(traqFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(serialize(entry))

	return err
}

func (fs *FileSystemStorage) Content(date time.Time) ([]string, error) {
	return fs.loader(fs.Path(date))
}
