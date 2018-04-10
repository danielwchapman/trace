// Package trace provides efficent and minimalist logging.
//
// Two logging levels are defined: trace and info. The trace level
// is for developers who debug code. The info level is for
// software operators (the folks running the code.) Examples of events
// the info level could include are logins, webpage loads,
// requests, hardware failures, or error events that cannot be
// handled gracefully.
//
// Logging groups are provided for organizing certain types of
// events and differientating their output location. For
// example, one might create an "Audit" group and output these logs
// to a file named "audit.log". There is only one logging group
// on initialization, the empty string. Additional groups can be
// defined with the RegisterGroup function.
//
// Logging groups and the entire trace level can be turned on or off
// depending on performance and requirements. For example, the trace
// level should typically be off in production systems. Groups and
// the trace level can be turned on or off while the software is running.
// The trace level is disabled by default.
package trace

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	// DefaultGroupId is the ID of the default logging group
	DefaultGroupId = 0

	// Number of logging requests and commands the channel buffer can hold
	chanBufSize = 1024
)

type level int

const (
	// Logging level for what developers care about
	trace level = iota + 1

	// Logging level for what software operators care about
	info
)

var (
	// Channel for ordering and concurrently outputing log messages
	logstream chan logApi

	// Tracks when logRoutine has completed all requests
	waitGroup sync.WaitGroup

	// Keeps all logging groups. Default group has index = 0 and name = ""
	groups []*groupData = make([]*groupData, 0, 4)

	// Indicates whether to output trace level logs
	traceEnabled bool = false
)

func init() {
	reset()
}

type groupData struct {
	name    string
	output  io.Writer
	enabled bool
}

type logApi interface {
	do()
}

type traceMsg struct {
	group int
	t     time.Time
	msg   string
}

func (m *traceMsg) do() {
	if traceEnabled && groups[m.group].enabled {
		printLog(m.group, m.t, m.msg)
	}
}

type infoMsg struct {
	group int
	t     time.Time
	msg   string
}

func (m *infoMsg) do() {
	if groups[m.group].enabled {
		printLog(m.group, m.t, m.msg)
	}
}

type cmdEnabletrace struct {
	on bool
}

func (c *cmdEnabletrace) do() {
	traceEnabled = c.on
}

type cmdEnableGroup struct {
	group int
	on    bool
}

func (c *cmdEnableGroup) do() {
	groups[c.group].enabled = c.on
}

// log is a helper function for processing new log requests from the caller
func log(group int, l level, format string, a ...interface{}) {
	t := time.Now()

	var m string
	if len(format) > 0 {
		m = fmt.Sprintf(format, a...)
	} else {
		m = fmt.Sprint(a...)
	}

	var cmd logApi
	if l == trace {
		cmd = &traceMsg{group: group, t: t, msg: m}
	} else if l == info {
		cmd = &infoMsg{group: group, t: t, msg: m}
	}

	logstream <- cmd
}

// logRoutine is a goroutine for outputing logging in parallel
func logRoutine() {
	for i := range logstream {
		i.do()
	}

	waitGroup.Done()
}

// printLog is a helper function for formating a log message
func printLog(group int, t time.Time, msg string) {
	strTime := t.UTC().Format("2006-1-2 15:04:05.000000")
	if group == DefaultGroupId {
		fmt.Fprintf(groups[DefaultGroupId].output, "%s %s\n", strTime, msg)
	} else {
		groupname := groups[group].name
		fmt.Fprintf(groups[group].output, "%s [%s] %s\n", strTime, groupname, msg)
	}
}

// reset is a helper function for initializing the trace package.
func reset() {
	if len(groups) == 0 {
		groups = append(groups, &groupData{name: "", output: os.Stdout, enabled: true})
	}

	logstream = make(chan logApi, chanBufSize)
	waitGroup.Add(1)
	go logRoutine()
}

// Done is called at end of program to ensure all logs are printed
func Done() {
	close(logstream)
	waitGroup.Wait()
}

// EnableGroup turns the group logging on or off
func EnableGroup(group int, on bool) {
	logstream <- &cmdEnableGroup{group, on}
}

// EnableTrace turns tracing level logging on or off
func EnableTrace(on bool) {
	logstream <- &cmdEnabletrace{on}
}

// Info logs a message to default group at info level. Similar to fmt.Print(...)
func Info(a ...interface{}) {
	log(0, info, "", a...)
}

// Infof logs a message to default group at info level. Similar to fmt.Printf(...)
func Infof(format string, a ...interface{}) {
	log(0, info, format, a...)
}

// Infog logs a message to given group at info level. Similar to fmt.Print(...)
func Infog(group int, a ...interface{}) {
	log(group, info, "", a...)
}

// Infogf logs a message to given group. Similar to fmt.Printf(...)
func Infogf(group int, format string, a ...interface{}) {
	log(group, info, format, a...)
}

// RegisterGroup registers a new logging group.
//
// It is to be called in a package's init() function. It returns a unique group ID
// for the calling package to store so it can later change the group configuration.
func RegisterGroup(name string, output io.Writer, on bool) int {
	for _, group := range groups {
		if name == group.name {
			panic("Group name already exists")
		}
	}

	if len(groups) == 0 {
		groups = append(groups, &groupData{name: "", output: os.Stdout, enabled: true})
	}

	groups = append(groups, &groupData{name: name, output: output, enabled: on})
	return len(groups) - 1
}

// SetDefaultOutput sets the output location of for the default logging group.
func SetDefaultOutput(output io.Writer) {
	if len(groups) == 0 {
		groups = append(groups, &groupData{name: "", output: output, enabled: true})
	} else {
		groups[0] = &groupData{name: "", output: output, enabled: groups[0].enabled}
	}
}

// Trace logs a message to default group at trace level. Similar to fmt.Print(...)
func Trace(a ...interface{}) {
	log(0, trace, "", a...)
}

// Trace logs a message to default group at trace level. Similar to fmt.Printf(...)
func Tracef(format string, a ...interface{}) {
	log(0, trace, format, a...)
}

// Traceg logs a message to given group at trace level. Similar to fmt.Print(...)
func Traceg(group int, a ...interface{}) {
	log(group, trace, "", a...)
}

// Tracegf logs a message to given group at trace level. Similar to fmt.Printf(...)
func Tracegf(group int, format string, a ...interface{}) {
	log(group, trace, format, a...)
}
