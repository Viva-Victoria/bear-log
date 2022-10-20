package bear_log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"os"
	"testing"
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

	logRandInt    = rand.Intn(2_000_000) - 1_000_000
	logRandString = randomText()
)

// BenchmarkBearLogger/instant
//BenchmarkBearLogger/instant-12          10330618              3741 ns/op
//              88 B/op          3 allocs/op
//BenchmarkBearLogger/format
//BenchmarkBearLogger/format-12            9008181              4225 ns/op
//             176 B/op          6 allocs/op
//BenchmarkBearLogger/fields
//BenchmarkBearLogger/fields-12             223656              4870 ns/op
//             921 B/op         14 allocs/op
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
			bearLog.DebugF("test %d %s", rand.Int(), logRandString)
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

// BenchmarkZapLog/instant
//BenchmarkZapLog/instant-12               9983971              3729 ns/op
//               0 B/op          0 allocs/op
//BenchmarkZapLog/format
//BenchmarkZapLog/format-12                8578303              4268 ns/op
//             102 B/op          3 allocs/op
//BenchmarkZapLog/fields
//BenchmarkZapLog/fields-12                8018008              4593 ns/op
//             128 B/op          1 allocs/op
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
			zapLogF.Debug("test %d %s", rand.Int(), logRandString)
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
