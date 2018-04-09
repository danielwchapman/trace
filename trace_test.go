package trace

import (
	"regexp"
	"testing"
)

const (
	timeFormat = `\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d\d:\d\d\.\d{6}`
)

// implements io.Writer
type memoryLog []string

func (l *memoryLog) Write(p []byte) (n int, err error) {
	*l = append(*l, string(p))
	return len(p), nil
}

func Test_Log(t *testing.T) {
	var logMemFile memoryLog
	logMemFile = make([]string, 0, 4)

	SetDefaultGroup(&logMemFile)

	Trace("Test trace")
	Info("Test info")

	msgNumber := 3
	anotherParam := "after others"
	Tracef("Test trace number %d %s", msgNumber, anotherParam)

	msgNumber = 4
	Infof("Test info number %d", msgNumber)

	Done()

	var gold []string
	gold = make([]string, 0, 4)
	gold = append(gold, timeFormat+` Test trace`)
	gold = append(gold, timeFormat+` Test info`)
	gold = append(gold, timeFormat+` Test trace number 3 after others`)
	gold = append(gold, timeFormat+` Test info number 4`)

	for i, line := range logMemFile {
		if match, err := regexp.MatchString(gold[i], line); err != nil || !match {
			t.Error("Trace failed: Line mismatch on line", i+1, "Recieved:\n", line)
		}
	}
}

func Test_LogGroup(t *testing.T) {
	reset()

	var logMemFile memoryLog
	logMemFile = make([]string, 0, 4)

	group := RegisterGroup("test", &logMemFile, true)

	Traceg(group, "Test trace")

	Infog(group, "Test info")

	msgNumber := 3
	anotherParam := "after others"
	Tracegf(group, "Test trace number %d %s", msgNumber, anotherParam)

	msgNumber = 4
	Infogf(group, "Test info number %d", msgNumber)

	Done()

	var gold []string
	gold = make([]string, 0, 4)
	gold = append(gold, `\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d\d:\d\d\.\d{6} \[test\] Test trace`)
	gold = append(gold, `\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d\d:\d\d\.\d{6} \[test\] Test info`)
	gold = append(gold, `\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d\d:\d\d\.\d{6} \[test\] Test trace number 3 after others`)
	gold = append(gold, `\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d\d:\d\d\.\d{6} \[test\] Test info number 4`)

	for i, line := range logMemFile {
		if match, err := regexp.MatchString(gold[i], line); err != nil || !match {
			t.Error("Trace failed: Line mismatch on line", i+1, "Recieved:\n", line)
		}
	}
}
