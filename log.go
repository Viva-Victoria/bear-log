package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger interface {
	WithTags(tags ...string) Logger
	WithFields(fields ...Field) Logger

	TraceEntry() Entry
	DebugEntry() Entry
	WarnEntry() Entry
	ErrorEntry() Entry
	CriticalEntry() Entry

	TraceF(format string, args ...any)
	DebugF(format string, args ...any)
	WarnF(format string, args ...any)
	ErrorF(format string, args ...any)
	CriticalF(format string, args ...any)

	Trace(message string)
	Debug(message string)
	Warn(message string)
	Error(message string)
	Critical(message string)
}

type FieldMapping struct {
	FormatLevel     LevelFormat
	TimeKey         string
	TimeFormat      string
	LevelKey        string
	MessageKey      string
	CallerKey       string
	StacktraceKey   string
	TagsKey         string
	FieldsNamespace string
}

type BearLogger struct {
	mutex   *sync.Mutex
	output  io.Writer
	format  FormatFunc
	now     func() time.Time
	mapping FieldMapping
	tags    []string
	fields  []Field
}

func NewBearLogger(options ...Option) BearLogger {
	initOptions := logOptions{
		output: os.Stdout,
		format: FormatJson,
		fieldMapping: FieldMapping{
			TimeKey:       "time",
			TimeFormat:    time.RFC3339,
			LevelKey:      "level",
			FormatLevel:   LevelFormatFunc(func(level Level) string { return level.String() }),
			MessageKey:    "message",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			TagsKey:       "tags",
		},
		tags:   nil,
		fields: nil,
	}
	for _, opt := range options {
		initOptions = opt.Apply(initOptions)
	}

	return BearLogger{
		mutex:   initOptions.mutex,
		output:  initOptions.output,
		format:  initOptions.format,
		mapping: initOptions.fieldMapping,
		tags:    initOptions.tags,
		fields:  initOptions.fields,
		now:     time.Now,
	}
}

func (b BearLogger) WithTags(tags ...string) Logger {
	if len(tags) == 0 {
		return b
	}

	b.tags = append(b.tags, tags...)
	return b
}

func (b BearLogger) WithFields(fields ...Field) Logger {
	if len(fields) == 0 {
		return b
	}

	b.fields = append(b.fields, fields...)
	return b
}

func (b BearLogger) TraceEntry() Entry {
	return b.entry(LevelTrace)
}

func (b BearLogger) DebugEntry() Entry {
	return b.entry(LevelDebug)
}

func (b BearLogger) WarnEntry() Entry {
	return b.entry(LevelWarn)
}

func (b BearLogger) ErrorEntry() Entry {
	return b.entry(LevelError)
}

func (b BearLogger) CriticalEntry() Entry {
	return b.entry(LevelCritical)
}

func (b BearLogger) TraceF(format string, args ...any) {
	b.writeF(LevelTrace, format, args...)
}

func (b BearLogger) Trace(message string) {
	b.write(LevelTrace, b.now(), message, nil, nil)
}

func (b BearLogger) DebugF(format string, args ...any) {
	b.writeF(LevelDebug, format, args...)
}

func (b BearLogger) Debug(message string) {
	b.write(LevelDebug, b.now(), message, nil, nil)
}

func (b BearLogger) WarnF(format string, args ...any) {
	b.writeF(LevelWarn, format, args...)
}

func (b BearLogger) Warn(message string) {
	b.write(LevelWarn, b.now(), message, nil, nil)
}

func (b BearLogger) ErrorF(format string, args ...any) {
	b.writeF(LevelError, format, args...)
}

func (b BearLogger) Error(message string) {
	b.write(LevelError, b.now(), message, nil, nil)
}

func (b BearLogger) CriticalF(format string, args ...any) {
	b.writeF(LevelCritical, format, args...)
}

func (b BearLogger) Critical(message string) {
	b.write(LevelCritical, b.now(), message, nil, nil)
}

func (b BearLogger) writeF(level Level, format string, args ...any) {
	b.write(level, time.Now(), fmt.Sprintf(format, args...), nil, nil)
}

func (b BearLogger) entry(level Level) Entry {
	return EntryImpl{
		level:     level,
		time:      time.Now(),
		writeFunc: b.write,
	}
}

func (b BearLogger) write(level Level, timestamp time.Time, message string, tags []string, fields []Field) {
	if b.mutex != nil {
		b.mutex.Lock()
		defer b.mutex.Unlock()
	}

	if len(tags) > 0 {
		if len(b.tags) > 0 {
			b.tags = append(b.tags, tags...)
		} else {
			b.tags = tags
		}
	}
	if len(fields) > 0 {
		if len(b.fields) > 0 {
			b.fields = append(b.fields, fields...)
		} else {
			b.fields = fields
		}
	}

	_, _ = b.output.Write(b.format(b.mapping, level, timestamp, message, b.tags, b.fields))
}
