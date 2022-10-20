package log

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestBearLogger(t *testing.T) {
	log := NewBearLogger(WithTags("t1"), WithFields(String("f1", "v1")))
	log.Debug("test")
}

func newZapLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.RFC3339TimeEncoder,
	}

	f, _ := os.OpenFile("zap.log", os.O_CREATE, os.ModePerm)
	return zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), f, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return true
	})), zap.Fields(zap.Array("tags", zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
		encoder.AppendString("t1")
		return nil
	})), zap.String("f1", "v1")))
}

func newBearLog() BearLogger {
	f, _ := os.OpenFile("bear.log", os.O_CREATE, os.ModePerm)
	return NewBearLogger(WithOutput(f), WithTags("t1"), WithFields(String("f1", "v1")))
}

var (
	zapLog  = newZapLogger()
	bearLog = newBearLog()

	logRandInt    = randomInt(-1_000_000, 1_000_000)
	logRandString = randomText()
)

func BenchmarkBearLogger(b *testing.B) {
	b.Run("instant", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			bearLog.Debug("test")
		}
	})
	b.Run("format", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			bearLog.DebugF("test %d %s", logRandInt, logRandString)
		}
	})
	b.Run("fields", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			bearLog.DebugEntry().
				Message("test").
				Tags("t3", "t4").
				Fields(String("f2", "v2")).
				Write()
		}
	})
}

func BenchmarkZapLog(b *testing.B) {
	b.Run("instant", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			zapLog.Debug("test")
		}
	})
	zapLogF := zapLog.Sugar()
	b.Run("format", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			zapLogF.Debug("test %d %s", logRandInt, logRandString)
		}
	})
	b.Run("fields", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			zapLog.Debug("test", zap.Array("tags", zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
				encoder.AppendString("t3")
				encoder.AppendString("t4")
				return nil
			})), zap.String("f2", "v2"))
		}
	})
}
