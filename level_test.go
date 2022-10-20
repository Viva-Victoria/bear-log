package bear_log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	_levelTestCases = map[Level]string{
		LevelTrace:    "TRACE",
		LevelDebug:    "DEBUG",
		LevelWarn:     "WARN",
		LevelError:    "ERROR",
		LevelCritical: "CRITICAL",
	}
)

func TestLevel_String(t *testing.T) {
	for level, text := range _levelTestCases {
		assert.Equal(t, text, level.String())
	}
	assert.Empty(t, Level(128).String())
}

func TestLevelEnablerFunc_IsEnabled(t *testing.T) {
	levelEnablerFunc := LevelEnablerFunc(func(level Level) bool {
		return level >= LevelError
	})

	for level := LevelTrace; level <= LevelCritical; level++ {
		assert.Equal(t, levelEnablerFunc(level), levelEnablerFunc.IsEnabled(level))
	}
}
