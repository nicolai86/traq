package main

import (
  "os"
  "fmt"
  "flag"
  "time"
  "io/ioutil"
  // "path/filepath"
)

var traqPath string = os.Getenv("TRAQ_DATA_DIR")
var month int
var year int
var day int
var project string = "timestamps"
var date string

func printFile(project string, date time.Time) {
  var traqFile = fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", traqPath, project, date.Year(), date.Year(), date.Month(), date.Day())
  var content, error = ioutil.ReadFile(traqFile)
  if (error == nil) {
    fmt.Print(string(content))
    fmt.Println("%%")
  } else {
    // fmt.Println(traqFile, " is unknown")
  }
}

func printMonth(project string, year int, month int) {
  var startDate time.Time = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
  for {
    printFile(project, startDate)
    startDate = startDate.Add(time.Hour * 24)
    if (int(startDate.Month()) != month) {
      break
    }
  }
}

func writeToFile(project string, date time.Time, command string) {
  var traqFile = fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", traqPath, project, date.Year(), date.Year(), date.Month(), date.Day())
  var file, error = os.OpenFile(traqFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
  if (error == nil) {
    var line = fmt.Sprintf("%s;%s;%s\n", date.Format("Mon Jan 2 15:04:05 -0700 2006"), command, "")
    file.WriteString(line)
    file.Close()
  }
}

func main() {
  flag.IntVar(&year      , "y", 0 , "print tracked times for a given year")
  flag.IntVar(&month     , "m", 0 , "print tracked times for a given month")

  flag.StringVar(&date   , "d", "" , "print tracked times for a given date")
  flag.StringVar(&project, "p", "" , "print data for a given project")

  flag.Parse()

  var now = time.Now()
  var t, error = time.Parse("2006-01-02", date)
  if error == nil {
    year = t.Year()
    month = int(t.Month())
    day = t.Day()
  } else {
    if (month == 0 && year == 0) {
      day = now.Day()
    } else {
      day = 1
    }
    if (year == 0) {
      year = now.Year()
    }
    if (month == 0) {
      month = int(now.Month())
    }
  }

  var command string = flag.Arg(0)
  if (command != "" && command != "stop") {
    command = "#" + command
  }

  if (command == "") {
    printMonth(project, year, month)
  } else {
    writeToFile(project, now, command)
  }
}