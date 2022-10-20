package log

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

type EntryImpl struct {
	time      time.Time
	writeFunc WriteFunc
	message   string
	tags      []string
	fields    []Field
	level     Level
}

func NewEntry(level Level, time time.Time, writeFunc WriteFunc) Entry {
	return EntryImpl{
		writeFunc: writeFunc,
		time:      time,
		level:     level,
	}
}

func (e EntryImpl) Message(m string) Entry {
	e.message = m
	return e
}

func (e EntryImpl) Format(f string, args ...any) Entry {
	return e.Message(fmt.Sprintf(f, args...))
}

func (e EntryImpl) Tags(tags ...string) Entry {
	e.tags = append(e.tags, tags...)
	return e
}

func (e EntryImpl) Fields(fields ...Field) Entry {
	e.fields = append(e.fields, fields...)
	return e
}

func (e EntryImpl) Write() {
	if e.writeFunc == nil {
		log.Println("writeFunc required!")
		return
	}

	e.writeFunc(e.level, e.time, e.message, e.tags, e.fields)
}
