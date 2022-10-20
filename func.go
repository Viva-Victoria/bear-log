package bear_log

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

	if len(fields) > 0 {
		if len(mapping.FieldsNamespace) > 0 {
			builder.Append(`,"`, mapping.FieldsNamespace, `":{`)
		} else {
			builder.Append(`,`)
		}

		for i, field := range fields {
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

		if len(mapping.FieldsNamespace) > 0 {
			builder.Append(`},`)
		}
	}

	return builder.Append(`}`, "\n").Bytes()
}
