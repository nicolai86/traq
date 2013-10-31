package traq

import (
  "fmt"
  "io/ioutil"
  "os"
  "time"
)

var traqPath string = os.Getenv("TRAQ_DATA_DIR")

func FilePath(project string, date time.Time) (path string) {
  return fmt.Sprintf("%s/%s/%d/%d-%02d-%02d", traqPath, project, date.Year(), date.Year(), date.Month(), date.Day())
}

func PrintFile(project string, date time.Time) {
  var content, error = ioutil.ReadFile(FilePath(project, date))

  if error == nil {
    fmt.Print(string(content))
    fmt.Println("%%")
  }
}

func PrintMonth(project string, year int, month int) {
  var startDate time.Time = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
  for {
    PrintFile(project, startDate)
    startDate = startDate.Add(time.Hour * 24)
    if int(startDate.Month()) != month {
      break
    }
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
