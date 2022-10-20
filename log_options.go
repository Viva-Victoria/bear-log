package log

import (
	"io"
	"sync"
)

type logOptions struct {
	mutex        *sync.Mutex
	output       io.Writer
	format       FormatFunc
	fieldMapping FieldMapping
	tags         []string
	fields       []Field
}

type Option interface {
	Apply(options logOptions) logOptions
}

type OptionFunc func(options logOptions) logOptions

// Apply
// nolint suppress revive
func (f OptionFunc) Apply(options logOptions) logOptions {
	return f(options)
}

func WithMutex(mx *sync.Mutex) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.mutex = mx
		return options
	})
}

func WithOutput(writer io.Writer) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.output = writer
		return options
	})
}

func WithFormat(format FormatFunc) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.format = format
		return options
	})
}

func WithFieldMapping(mapping FieldMapping) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping = mapping
		return options
	})
}

func WithTags(tags ...string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.tags = append(options.tags, tags...)
		return options
	})
}

func WithFields(fields ...Field) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fields = append(options.fields, fields...)
		return options
	})
}

func WithTime(key, format string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.TimeKey = key
		options.fieldMapping.TimeFormat = format
		return options
	})
}

func WithLevel(key string, format LevelFormat) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.LevelKey = key
		options.fieldMapping.FormatLevel = format
		return options
	})
}

func WithMessageKey(key string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.MessageKey = key
		return options
	})
}

func WithCallerKey(key string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.CallerKey = key
		return options
	})
}

func WithStacktraceKey(key string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.StacktraceKey = key
		return options
	})
}

func WithTagsKey(key string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.TagsKey = key
		return options
	})
}

func WithFieldsNamespace(key string) Option {
	return OptionFunc(func(options logOptions) logOptions {
		options.fieldMapping.FieldsNamespace = key
		return options
	})
}
