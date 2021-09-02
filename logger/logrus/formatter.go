// Package logrus
//
// @author: xwc1125
// @date: 2021/8/30
package logrus

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

var (
	termTimeFormat = "2006-01-02 15:04:05.000"
	moduleKey      = "module"
	termMsgJust    = 40
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// UseColor - default: disable colors
	UseColor bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// LocationEnabled - print caller
	LocationEnabled bool
	LocationTrims   []string
	locationLength  uint32

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string

	Modules    []string
	modulesReg atomic.Value
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = termTimeFormat
	}
	module := ToString(entry.Data[moduleKey])
	if !f.isPrint(module) {
		return []byte{}, nil
	}

	// output buffer
	b := &bytes.Buffer{}

	// write level
	var level = strings.ToUpper(entry.Level.String())
	{
		levelLen := len(level)
		if f.UseColor && levelColor > 0 {
			level = colorRender(level, levelColor, 0)
		}
		b.WriteString(level)
		b.WriteString(strings.Repeat(" ", 5-levelLen))
	}

	// write time
	{
		b.WriteString(" ")
		b.WriteString("[")
		b.WriteString(entry.Time.Format(timestampFormat))
		b.WriteString("]")
	}

	// write module
	{
		if levelColor > 0 {
			module = colorRender(module, levelColor, 0)
			b.WriteString(" ")
			b.WriteString("[")
			b.WriteString(module)
			b.WriteString("]")
		}
	}

	// write location
	//var location string
	{
		//location := fmt.Sprintf("%+v", r.Call)
		//if f.LocationTrims!=nil&& len(f.LocationTrims)>0 {
		//	for _, prefix := range f.LocationTrims {
		//
		//	}
		//}
		if f.LocationEnabled {
			f.writeCaller(b, entry)
		}
	}

	// write message
	if f.TrimMessages {
		b.WriteString(" ")
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(" ")
		b.WriteString(entry.Message)
	}

	// write padding
	{
		// try to justify the log output for short messages
		length := utf8.RuneCountInString(entry.Message)
		justLen := termMsgJust - utf8.RuneCountInString(module)
		if len(entry.Data) > 0 {
			if length < justLen {
				b.Write(bytes.Repeat([]byte{' '}, justLen-length))
			} else {
				b.WriteString(" --> ")
			}
		}
	}

	// write fields
	{
		b.WriteString(" ")
		f.writeFields(b, entry, levelColor)
	}
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func colorRender(str string, color int, weight int, extraArgs ...interface{}) string {
	//闪烁效果
	isBlink := int64(0)
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	//下划线效果
	isUnderLine := int64(0)
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" (%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry, levelColor int) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field, levelColor)
			b.WriteString(" ")
		}
	}
}

//func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
//	length := len(entry.Data)
//	foundFieldsMap := map[string]bool{}
//	for _, field := range f.FieldsOrder {
//		if _, ok := entry.Data[field]; ok {
//			foundFieldsMap[field] = true
//			length--
//			f.writeField(b, entry, field)
//		}
//	}
//
//	if length > 0 {
//		notFoundFields := make([]string, 0, length)
//		for field := range entry.Data {
//			if foundFieldsMap[field] == false {
//				notFoundFields = append(notFoundFields, field)
//			}
//		}
//
//		sort.Strings(notFoundFields)
//
//		for _, field := range notFoundFields {
//			f.writeField(b, entry, field)
//		}
//	}
//}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string, levelColor int) {
	if !strings.EqualFold(field, moduleKey) {
		fmt.Fprintf(b, "%s", colorRender(field, levelColor, 0))
		fmt.Fprintf(b, "=%v", entry.Data[field])
	}

	//if !f.NoFieldsSpace {
	//	b.WriteString(" ")
	//}
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.FatalLevel, logrus.PanicLevel:
		return 35
	case logrus.ErrorLevel:
		return 31
	case logrus.WarnLevel:
		return 33
	case logrus.InfoLevel:
		return 32
	case logrus.DebugLevel:
		return 36
	case logrus.TraceLevel:
		return 34
	default:
		return 0
	}
}

func (f *Formatter) isPrint(module string) bool {
	if f.Modules == nil || len(f.Modules) == 0 {
		return true
	}
	if len(strings.TrimSpace(module)) == 0 {
		return false
	}
	var modulesRegMap sync.Map
	if f.modulesReg.Load() != nil {
		modulesRegMap = f.modulesReg.Load().(sync.Map)
	} else {
		for _, s := range f.Modules {
			modulesRegMap.Store(s, s)
		}
		f.modulesReg.Store(modulesRegMap)
	}
	if _, ok := modulesRegMap.Load(module); ok {
		return true
	}
	return false
}
