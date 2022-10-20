package log

import (
	"bytes"
	crypto_rand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	math_rand "math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type NewEntryTestCase struct {
	writeFunc func() (WriteFunc, *bool)
	time      time.Time
	level     Level
}

var (
	_nilWriteFunc      = func() (WriteFunc, *bool) { return nil, nil }
	_newEntryTestCases = map[string]NewEntryTestCase{
		"debug/utc/nil": {
			level:     LevelDebug,
			time:      time.Date(2022, 1, 1, 10, 45, 0, 0, time.UTC),
			writeFunc: _nilWriteFunc,
		},
		"error/local-0/nil": {
			level:     LevelError,
			time:      time.Date(0, 1, 1, 0, 0, 0, 0, time.Local),
			writeFunc: _nilWriteFunc,
		},
		"custom/local-2020/func": {
			level: Level(128),
			time:  time.Date(2020, 5, 17, 15, 27, 18, 0, time.Local),
			writeFunc: func() (WriteFunc, *bool) {
				value := false
				f := &value
				return func(Level, time.Time, string, []string, []Field) {
					*f = true
				}, f
			}},
	}
)

func TestNewEntry(t *testing.T) {
	for name, testCase := range _newEntryTestCases {
		t.Run(name, func(t *testing.T) {
			writeFunc, f := testCase.writeFunc()

			entry, ok := NewEntry(testCase.level, testCase.time, writeFunc).(EntryImpl)
			require.True(t, ok)

			assert.Equal(t, testCase.level, entry.level)
			assert.Equal(t, testCase.time, entry.time)

			if writeFunc == nil {
				assert.Nil(t, entry.writeFunc)
			} else {
				require.NotNil(t, entry.writeFunc)
				entry.writeFunc(entry.level, entry.time, entry.message, entry.tags, entry.fields)

				require.NotNil(t, f)
				assert.True(t, *f)
			}
		})
	}
}

func TestLogEntry_Message(t *testing.T) {
	var entry EntryImpl

	assert.Empty(t, entry.message)

	text := randomText()
	newEntry := entry.Message(text)
	assert.Empty(t, entry.message)

	entry, ok := newEntry.(EntryImpl)
	require.True(t, ok)
	assert.Equal(t, text, entry.message)
}

func TestLogEntry_Format(t *testing.T) {
	var entry EntryImpl

	assert.Empty(t, entry.message)

	var (
		text   = randomText()
		number = randomInt(-1_000_000, 1_000_000)
		format = "%s %d"
	)
	newEntry := entry.Format(format, text, number)
	assert.Empty(t, entry.message)

	entry, ok := newEntry.(EntryImpl)
	require.True(t, ok)
	assert.Equal(t, fmt.Sprintf(format, text, number), entry.message)
}

func TestLogEntry_Tags(t *testing.T) {
	var entry EntryImpl
	assert.Empty(t, entry.tags)

	tags := randomTags()
	newEntry := entry.Tags(tags...)
	assert.Empty(t, entry.tags)

	entry, ok := newEntry.(EntryImpl)
	require.True(t, ok)

	assert.Equal(t, tags, entry.tags)

	tags2 := randomTags()
	entry, ok = entry.Tags(tags2...).(EntryImpl)
	require.True(t, ok)

	assert.Equal(t, append(tags, tags2...), entry.tags)
}

func TestLogEntry_Fields(t *testing.T) {
	var entry EntryImpl
	assert.Empty(t, entry.tags)

	fields := randomFields()
	newEntry := entry.Fields(fields...)
	assert.Empty(t, entry.fields)

	entry, ok := newEntry.(EntryImpl)
	require.True(t, ok)

	assert.Equal(t, fields, entry.fields)

	fields2 := randomFields()
	entry, ok = entry.Fields(fields2...).(EntryImpl)
	require.True(t, ok)

	assert.Equal(t, append(fields, fields2...), entry.fields)
}

func TestLogEntry_Write(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var entry EntryImpl

		var buffer bytes.Buffer
		log.SetFlags(0)
		log.SetOutput(&buffer)
		entry.Write()
		log.SetOutput(os.Stdout)
		log.SetFlags(log.LstdFlags)

		assert.Equal(t, "writeFunc required!\n", buffer.String())
	})
	t.Run("call", func(t *testing.T) {
		var buffer bytes.Buffer

		entry := EntryImpl{
			level:   LevelDebug,
			time:    time.Now(),
			message: randomText(),
			writeFunc: func(level Level, timestamp time.Time, message string, tags []string, fields []Field) {
				buffer.WriteString(fmt.Sprintf("%v %s %s %s %s", level, timestamp, message, tags, fields))
			},
		}

		entry.Write()

		assert.Equal(t, fmt.Sprintf("%v %s %s %s %s", entry.level, entry.time, entry.message, entry.tags, entry.fields), buffer.String())
	})
}

var randomText = uuid.NewString

func randomInt(min, max int) int {
	bigInt, _ := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(max-min)))
	return int(bigInt.Int64()) + min
}

func randomUInt() uint {
	return uint(randomInt(0, 4294967295))
}

func randomFloat() float64 {
	return math_rand.Float64()
}

func randomTags() []string {
	count := randomInt(5, 10)
	array := make([]string, 0, count)
	for i := 0; i < cap(array); i++ {
		array = append(array, randomText())
	}
	return array
}

func randomFields() []Field {
	count := randomInt(5, 10)
	array := make([]Field, 0, count)
	for i := 0; i < cap(array); i++ {
		if randomFloat() > 0.5 {
			array = append(array, String(randomText(), randomText()))
		} else {
			array = append(array, Int(randomText(), randomInt(0, 10_000_000)))
		}
	}
	return array
}
