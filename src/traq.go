package traq

import (
  "fmt"
  "io/ioutil"
  "os"
  "time"
  "strings"
)

var traqPath string = os.Getenv("TRAQ_DATA_DIR")

func FilePath(project string, date time.Time) (path string) {
  return fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", traqPath, project, date.Year(), date.Year(), date.Month(), date.Day())
}

func DatesInMonth(year int, month int) ([]time.Time) {
  var dates []time.Time = make([]time.Time, 0)
  var date time.Time = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

  for {
    dates = append(dates, date)
    date = date.Add(time.Hour * 24)
    if int(date.Month()) != month {
      break
    }
  }

  return dates
}

func SumFile(content string) (map[string]int64, error) {
  var totalled map[string]int64 = make(map[string]int64)

  var lines []string = strings.Split(content, "\n")
  var currentTag string = ""
  var currentTime time.Time

  for _, line := range lines {
    if line != "" {
      var parts []string = strings.Split(line, ";")

      var t, error = time.Parse("Mon Jan 2 15:04:05 -0700 2006", parts[0])
      if error == nil {
        if parts[1] == "" {
          currentTag = parts[1]
          totalled[currentTag] = 0
        } else if parts[1] == "stop" {
          var diff = t.Unix() - currentTime.Unix()
          totalled[currentTag] = totalled[currentTag] + diff
          currentTag = ""
        } else if currentTag != parts[1] {
          var diff = t.Unix() - currentTime.Unix()
          totalled[currentTag] = totalled[currentTag] + diff
          currentTag = parts[1]
        }

        currentTime = t
      } else {
        return totalled, error
      }
    }
  }

  delete(totalled, "")

  return totalled, nil
}

func PrintDate(project string, date time.Time) {
  var content, error = ioutil.ReadFile(FilePath(project, date))

  if error == nil {
    fmt.Print(string(content))
    fmt.Println("%%")
  }
}

func EvaluateDate(project string, date time.Time) {
  var content, error = ioutil.ReadFile(FilePath(project, date))

  if error == nil {
    fmt.Printf("%d-%02d-%02d\n", date.Year(), date.Month(), date.Day())
    var totalled, _ = SumFile(string(content))
    // TODO handle errors
    for key, value := range totalled {
      fmt.Printf("%s:%2.4f\n", key, float64(value) / 60.0 / 60.0)
    }

    fmt.Println("%%")
  }
}

func PrintMonth(project string, year int, month int) {
  for _, date := range DatesInMonth(year, month) {
    PrintDate(project, date)
  }
}

func EvaluateMonth(project string, year int, month int) {
  for _, date := range DatesInMonth(year, month) {
    EvaluateDate(project, date)
  }
}

func WriteToFile(project string, date time.Time, command string) {
  var file, error = os.OpenFile(FilePath(project, date), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
  if error == nil {
    var line = fmt.Sprintf("%s;%s;%s\n", date.Format("Mon Jan 2 15:04:05 -0700 2006"), command, "")
    file.WriteString(line)
    file.Close()
  }
}