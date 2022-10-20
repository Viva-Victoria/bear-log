package bear_log

import (
	"fmt"
	"log"
	"time"
)

type Entry interface {
	Message(string) Entry
	Format(string, ...any) Entry
	Tags(...string) Entry
	Fields(...Field) Entry
	Write()
}

type LogEntry struct {
	tags      []string
	fields    []Field
	time      time.Time
	message   string
	level     Level
	writeFunc WriteFunc
}

func NewEntry(level Level, time time.Time, writeFunc WriteFunc) Entry {
	return LogEntry{
		writeFunc: writeFunc,
		time:      time,
		level:     level,
	}
}

func (e LogEntry) Message(m string) Entry {
	e.message = m
	return e
}

func (e LogEntry) Format(f string, args ...any) Entry {
	return e.Message(fmt.Sprintf(f, args...))
}

func (e LogEntry) Tags(tags ...string) Entry {
	e.tags = append(e.tags, tags...)
	return e
}

func (e LogEntry) Fields(fields ...Field) Entry {
	e.fields = append(e.fields, fields...)
	return e
}

func (e LogEntry) Write() {
	if e.writeFunc == nil {
		log.Println("writeFunc required!")
		return
	}

	e.writeFunc(e.level, e.time, e.message, e.tags, e.fields)
}
