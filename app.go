package main

import (
  "flag"
  "time"
  "traq"
)

var month int
var year int
var day int
var project string = "timestamps"
var date string
var evaluate bool

func main() {
  flag.BoolVar(&evaluate, "e", false, "evaluate tracked times")
  flag.IntVar(&year, "y", 0, "print tracked times for a given year")
  flag.IntVar(&month, "m", 0, "print tracked times for a given month")

  flag.StringVar(&date, "d", "", "print tracked times for a given date")
  flag.StringVar(&project, "p", "", "print data for a given project")

  flag.Parse()

  var now = time.Now()
  var t, error = time.Parse("2006-01-02", date)
  if error == nil {
    year = t.Year()
    month = int(t.Month())
    day = t.Day()
  } else {
    if month == 0 && year == 0 {
      day = now.Day()
    } else {
      day = 1
    }
    if year == 0 {
      year = now.Year()
    }
    if month == 0 {
      month = int(now.Month())
    }
  }

  if evaluate {
    if date == "" {
      traq.EvaluateMonth(project, year, month)
    } else {
      traq.EvaluateDate(project, t)
    }
    return
  }

  var command string = flag.Arg(0)
  if command != "" && command != "stop" {
    command = "#" + command
  }

  if command == "" {
    if date == "" {
      traq.PrintMonth(project, year, month)
    } else {
      traq.PrintDate(project, t)
    }
  } else {
    traq.WriteToFile(project, now, command)
  }
}
