package bear_log

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"os"
	"testing"
	"time"
)

type formatJsonEntry struct {
	Level       Level       `json:"-"`
	LevelString string      `json:"level"`
	Time        time.Time   `json:"-"`
	TimeString  string      `json:"time"`
	Message     string      `json:"message"`
	Tags        []string    `json:"tags,omitempty"`
	Fields      fieldsArray `json:"fields,omitempty"`
}

type fieldsArray []Field

func (f fieldsArray) MarshalJSON() ([]byte, error) {
	result := make(map[string]any, len(f))

	for _, l := range f {
		var value any
		switch l.Type() {
		case TypeInt:
			value = l.Int()

		case TypeUInt:
			value = l.UInt()

		case TypeFloat:
			value = l.Float()

		case TypeBinary:
			array, ok := l.Addressable().([]byte)
			if !ok {
				return nil, fmt.Errorf("%v is not []byte", l.Addressable())
			}

			value = base64.StdEncoding.EncodeToString(array)

		case TypeString:
			value = l.StringValue()

		case TypeArray, TypeMap, TypeAny:
			data, err := json.Marshal(l.Addressable())
			if err != nil {
				return nil, err
			}

			value = string(data)
		}

		if errF, ok := l.(errorField); ok {
			_, err := errF.String()
			result[l.Key()] = fmt.Sprintf("%%(!%v)", err)
		} else {
			result[l.Key()] = value
		}
	}

	return json.Marshal(result)
}

type formatJsonTestCase struct {
	name    string
	entry   formatJsonEntry
	mapping FieldMapping
}

type errorField struct {
	LogField
}

func (e errorField) String() (string, error) {
	return "", ErrTest
}

var (
	ErrTest              = errors.New("test error")
	_defaultFieldMapping = FieldMapping{
		TimeKey:    "time",
		TimeFormat: time.RFC3339,
		LevelKey:   "level",
		FormatLevel: LevelFormatFunc(func(level Level) string {
			return level.String()
		}),
		MessageKey:      "message",
		CallerKey:       "caller",
		StacktraceKey:   "stacktrace",
		TagsKey:         "tags",
		FieldsNamespace: "fields",
	}

	_formatJsonTestCases = []formatJsonTestCase{
		generateFormatJsonTestCase(_defaultFieldMapping, nil, nil),
		generateFormatJsonTestCase(_defaultFieldMapping, []string{randomText(), randomText()}, nil),
		generateFormatJsonTestCase(_defaultFieldMapping, nil, []Field{
			Int(randomText(), rand.Int()),
			Float(randomText(), rand.Float64()),
			String(randomText(), randomText()),
			errorField{LogField{
				key:       randomText(),
				valueType: TypeString,
			}},
		}),
	}
	_benchmarkTestCase = generateFormatJsonTestCase(_defaultFieldMapping, []string{randomText(), randomText(), randomText()}, []Field{
		String(randomText(), randomText()),
		Int(randomText(), rand.Int()),
		Uint(randomText(), rand.Uint64()),
		Float(randomText(), rand.Float64()),
	})
)

func TestFormatJson(t *testing.T) {
	for _, testCase := range _formatJsonTestCases {
		entry := testCase.entry
		actual := FormatJson(testCase.mapping, entry.Level, entry.Time, entry.Message, entry.Tags, entry.Fields)

		expected, err := json.Marshal(entry)
		require.NoError(t, err)

		assert.JSONEq(t, string(expected), string(actual))
	}

	testCase := generateFormatJsonTestCase(FieldMapping{
		TimeKey:       "time",
		TimeFormat:    time.RFC3339,
		LevelKey:      "level",
		FormatLevel:   LevelFormatFunc(func(level Level) string { return level.String() }),
		MessageKey:    "message",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		TagsKey:       "tags",
	}, []string{randomText(), randomText()}, []Field{String(randomText(), randomText()), Int(randomText(), rand.Int())})
	result := map[string]any{
		testCase.mapping.LevelKey:   testCase.entry.Level.String(),
		testCase.mapping.TimeKey:    testCase.entry.Time.Format(testCase.mapping.TimeFormat),
		testCase.mapping.MessageKey: testCase.entry.Message,
		testCase.mapping.TagsKey:    testCase.entry.Tags,
	}
	for _, field := range testCase.entry.Fields {
		value, err := field.Value()
		if err != nil {
			require.NoError(t, err)
		}

		result[field.Key()] = value
	}

	expected, err := json.Marshal(result)
	require.NoError(t, err)

	actual := FormatJson(testCase.mapping, testCase.entry.Level, testCase.entry.Time, testCase.entry.Message, testCase.entry.Tags, testCase.entry.Fields)
	assert.JSONEq(t, string(expected), string(actual))
}

func BenchmarkFormatJson(b *testing.B) {
	b.ReportAllocs()

	var buffer bytes.Buffer

	for i := 0; i < b.N; i++ {
		entry := _benchmarkTestCase.entry
		actual := FormatJson(_benchmarkTestCase.mapping, entry.Level, entry.Time, entry.Message, entry.Tags, entry.Fields)
		buffer.Write(actual)
	}

	b.StopTimer()
	writeToFile("format", &buffer)
}

func BenchmarkJsonMarshal(b *testing.B) {
	entry := _benchmarkTestCase.entry

	result := map[string]any{
		"level":   entry.Level.String(),
		"time":    entry.Time.Format(time.RFC3339),
		"message": entry.Message,
		"tags":    entry.Tags,
	}

	if len(entry.Fields) > 0 {
		fields := make(map[string]any)
		for _, f := range entry.Fields {
			value, err := f.Value()
			require.NoError(b, err)

			fields[f.Key()] = value
		}

		result["fields"] = fields
	}

	b.ReportAllocs()
	b.ResetTimer()

	var buffer bytes.Buffer

	for i := 0; i < b.N; i++ {

		actual, _ := json.Marshal(entry)
		buffer.Write(actual)
	}

	b.StopTimer()
	writeToFile("marshal", &buffer)
}

func BenchmarkZap(b *testing.B) {
	b.ReportAllocs()

	entry := _benchmarkTestCase.entry
	fields := []zap.Field{zap.Array("tags", zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
		for _, tag := range entry.Tags {
			encoder.AppendString(tag)
		}
		return nil
	}))}

	if len(entry.Fields) > 0 {
		fieldsMap := make(map[string]any)
		for _, f := range entry.Fields {
			value, err := f.Value()
			require.NoError(b, err)

			fieldsMap[f.Key()] = value
		}

		fields = append(fields, zap.Any("fields", fieldsMap))
	}

	b.ResetTimer()

	var buffer bytes.Buffer

	for i := 0; i < b.N; i++ {
		entry := _benchmarkTestCase.entry

		encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
		})
		buf, _ := encoder.EncodeEntry(zapcore.Entry{
			Level:   zap.DebugLevel,
			Time:    entry.Time,
			Message: entry.Message,
		}, fields)
		buffer.Write(buf.Bytes())
	}

	b.StopTimer()
	writeToFile("zap", &buffer)
}

func generateFormatJsonTestCase(fm FieldMapping, tags []string, fields []Field) formatJsonTestCase {
	t := time.Now()
	message := randomText()

	return formatJsonTestCase{
		name: "no tags, no fields",
		entry: formatJsonEntry{
			Level:       LevelDebug,
			LevelString: LevelDebug.String(),
			Time:        t,
			TimeString:  t.Format(time.RFC3339),
			Message:     message,
			Tags:        tags,
			Fields:      fields,
		},
		mapping: fm,
	}
}

func writeToFile(mode string, b *bytes.Buffer) {
	f, _ := os.OpenFile(fmt.Sprintf("out/%s_%s", mode, randomText()), os.O_CREATE, os.ModePerm)
	_, _ = b.WriteTo(f)
}
