package logger

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func Initialize(colors bool, level logrus.Level) {
	formatter := new(log_formatter)
	formatter.CustomCallerFormatter = func(f *runtime.Frame) string {
		return fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
	}
	formatter.ShowFullLevel = true
	formatter.NoFieldsSpace = true
	formatter.MessageSeparator = "-"
	formatter.TimestampFormat = "2006-01-02T15:04:05.000000-0700"
	formatter.NoTimestampColor = false
	formatter.NoCallerColor = false
	formatter.NoFieldsColors = false
	formatter.NoColors = !colors
	formatter.HideKeys = true

	logrus.SetFormatter(formatter)
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)

	logrus.Infof(LOGGER_CONFIG, colors, level)
	logrus.Infof(MAIN_LOGICAL_CORES, runtime.NumCPU())
}

// log_formatter - logrus formatter, implements logrus.log_formatter
type log_formatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - do not apply color on fields
	NoFieldsColors bool

	// NoTimestampColor - do not apply color on timestamp
	NoTimestampColor bool

	// NoCallerColor - do not apply color on caller
	NoCallerColor bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string

	// MessageSeparator - message separator string
	MessageSeparator string
}

// Format an log entry
func (f *log_formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Output buffer
	b := &bytes.Buffer{}

	// Write timestamp
	if !f.NoColors && !f.NoTimestampColor {
		set_color(b, get_timestamp_color())
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMicro
	}
	b.WriteString(entry.Time.Format(timestampFormat))

	// Write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if !f.NoColors {
		set_color(b, get_level_color(entry.Level))
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	// Write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	// Write caller information
	if !f.NoColors && !f.NoCallerColor {
		set_color(b, get_caller_color())
	}

	f.writeCaller(b, entry)

	// write message
	if !f.NoColors {
		set_color(b, get_message_color())
	}

	if f.MessageSeparator == "" {
		b.WriteString(entry.Message)
	} else {
		b.WriteString(f.MessageSeparator + " " + entry.Message)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *log_formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprint(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				"(%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
		fmt.Fprintf(b, " ")
	}
}

func (f *log_formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *log_formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if !foundFieldsMap[field] {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *log_formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

const (
	colorRed     = 31
	colorGreen   = 32
	colorYellow  = 33
	colorBlue    = 34
	colorMagenta = 35
	colorCyan    = 36
	colorWhite   = 37
)

func get_level_color(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return colorWhite
	case logrus.WarnLevel, logrus.DebugLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

func get_timestamp_color() int {
	return colorMagenta
}

func get_caller_color() int {
	return colorGreen
}

func get_message_color() int {
	return colorWhite
}

func set_color(b *bytes.Buffer, color int) {
	fmt.Fprintf(b, "\x1b[%dm", color)
}
