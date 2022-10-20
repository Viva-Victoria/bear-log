package bear_log

type Level uint8

const (
	LevelTrace Level = iota + 1
	LevelDebug
	LevelWarn
	LevelError
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelCritical:
		return "CRITICAL"
	default:
		return ""
	}
}

type LevelEnabler interface {
	IsEnabled(level Level) bool
}

type LevelEnablerFunc func(level Level) bool

func (f LevelEnablerFunc) IsEnabled(level Level) bool {
	return f(level)
}
