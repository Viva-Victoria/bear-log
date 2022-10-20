package log

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptionFunc_Apply(t *testing.T) {
	var called bool
	f := OptionFunc(func(options logOptions) logOptions {
		called = true
		return options
	})

	f.Apply(logOptions{})
	assert.True(t, called)
}

func TestWithMutex(t *testing.T) {
	mx := &sync.Mutex{}

	opt := WithMutex(mx).Apply(logOptions{})
	assert.Equal(t, mx, opt.mutex)
}

func TestWithOutput(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	opt := WithOutput(buf).Apply(logOptions{})
	assert.Equal(t, buf, opt.output)
}

func TestWithFormat(t *testing.T) {
	data := []byte(randomText())
	testFormat := func(mapping FieldMapping, level Level, timestamp time.Time, message string, tags []string, fields []Field) []byte {
		return data
	}

	opt := WithFormat(testFormat).Apply(logOptions{})
	assert.Equal(t, data, opt.format(FieldMapping{}, LevelDebug, time.Now(), "", nil, nil))
}

func TestWithFieldMapping(t *testing.T) {
	var fm FieldMapping

	opt := WithFieldMapping(fm).Apply(logOptions{})
	assert.Equal(t, fm, opt.fieldMapping)
}

func TestWithTags(t *testing.T) {
	tags := []string{randomText(), randomText()}

	opt := WithTags(tags...).Apply(logOptions{})
	assert.Equal(t, tags, opt.tags)
}

func TestWithFields(t *testing.T) {
	fields := []Field{
		String(randomText(), randomText()),
		Int(randomText(), randomInt(-1000, 1000)),
	}

	opt := WithFields(fields...).Apply(logOptions{})
	assert.Equal(t, fields, opt.fields)
}

func TestWithTime(t *testing.T) {
	key := randomText()
	format := "2006-01-02"

	opt := WithTime(key, format).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.TimeKey)
	assert.Equal(t, format, opt.fieldMapping.TimeFormat)
}

func TestWithLevel(t *testing.T) {
	key := randomText()
	format := LevelFormatFunc(func(level Level) string {
		return level.String()
	})

	opt := WithLevel(key, format).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.LevelKey)
	assert.Equal(t, format(LevelCritical), opt.fieldMapping.FormatLevel.Format(LevelCritical))
}

func TestWithMessageKey(t *testing.T) {
	key := randomText()

	opt := WithMessageKey(key).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.MessageKey)
}

func TestWithCallerKey(t *testing.T) {
	key := randomText()

	opt := WithCallerKey(key).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.CallerKey)
}

func TestWithStacktraceKey(t *testing.T) {
	key := randomText()

	opt := WithStacktraceKey(key).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.StacktraceKey)
}

func TestWithTagsKey(t *testing.T) {
	key := randomText()

	opt := WithTagsKey(key).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.TagsKey)
}

func TestWithFieldsNamespace(t *testing.T) {
	key := randomText()

	opt := WithFieldsNamespace(key).Apply(logOptions{})
	assert.Equal(t, key, opt.fieldMapping.FieldsNamespace)
}
