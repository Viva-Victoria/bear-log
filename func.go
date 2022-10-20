package log

import (
	"time"
)

type LevelFormat interface {
	Format(level Level) string
}

type LevelFormatFunc func(level Level) string

func (f LevelFormatFunc) Format(level Level) string {
	return f(level)
}

type WriteFunc func(level Level, timestamp time.Time, message string, tags []string, fields []Field)

type FormatFunc func(mapping FieldMapping, level Level, timestamp time.Time, message string, tags []string, fields []Field) []byte

func FormatJson(mapping FieldMapping, level Level, timestamp time.Time, message string, tags []string, fields []Field) []byte {
	builder := NewStringBuilder()
	defer builder.Dispose()

	builder.Append(`{"`, mapping.LevelKey, `":"`, mapping.FormatLevel.Format(level), `","`, mapping.TimeKey, `":"`, timestamp.Format(mapping.TimeFormat), `"`)
	if len(message) > 0 {
		builder.Append(`,"`, mapping.MessageKey, `":"`, message, `"`)
	}
	formatJsonTags(mapping, tags, builder)
	formatJsonFields(mapping, fields, builder)

	return builder.Append(`}`, "\n").Bytes()
}

func formatJsonFields(mapping FieldMapping, fields []Field, builder StringBuilder) {
	if len(fields) > 0 {
		if len(mapping.FieldsNamespace) > 0 {
			builder.Append(`,"`, mapping.FieldsNamespace, `":{`)
		} else {
			builder.Append(`,`)
		}

		for i, field := range fields {
			formatJsonField(field, builder, i, fields)
		}

		if len(mapping.FieldsNamespace) > 0 {
			builder.Append(`},`)
		}
	}
}

func formatJsonField(field Field, builder StringBuilder, i int, fields []Field) {
	value, err := field.String()
	if err == nil {
		if field.Type() == TypeString || field.Type() == TypeBinary {
			builder.Append(`"`, value, `"`)
		} else {
			builder.Append(value)
		}
	}

	if i < len(fields)-1 {
		builder.Append(`,`)
	}
}

func formatJsonTags(mapping FieldMapping, tags []string, builder StringBuilder) {
	if len(tags) > 0 {
		builder.Append(`,"`, mapping.TagsKey, `":[`)
		for i, tag := range tags {
			builder.Append(`"`, tag, `"`)
			if i < len(tags)-1 {
				builder.Append(`,`)
			}
		}
		builder.Append(`]`)
	}
}
