package zapjournald

import (
	"testing"

	"github.com/ssgreg/journald"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type nopSync struct{}

func (nopSync) Write([]byte) (int, error) { return 0, nil }
func (nopSync) Sync() error               { return nil }

func BenchmarkLogger(b *testing.B) {
	zc := zap.NewProductionConfig()
	zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zc.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	b.Run("standard", func(b *testing.B) {
		encoder := zapcore.NewJSONEncoder(zc.EncoderConfig)
		core := zapcore.NewCore(encoder, nopSync{}, zc.Level)
		coreWithContext := core.With([]zapcore.Field{
			SyslogFacility(LogDaemon),
			SyslogIdentifier(),
			SyslogPid()})
		l := zap.New(coreWithContext)
		benchmarkLog(b, l)
	})
	b.Run("journald", func(b *testing.B) {
		encoder := zapcore.NewJSONEncoder(zc.EncoderConfig)
		core := NewCore(zc.Level, encoder, &journald.Journal{}, SyslogFields)
		core.j.TestModeEnabled = true // Disable actual writing to the journal.

		coreWithContext := core.With([]zapcore.Field{
			SyslogFacility(LogDaemon),
			SyslogIdentifier(),
			SyslogPid(),
		})
		l := zap.New(coreWithContext)
		benchmarkLog(b, l)
	})
}

func benchmarkLog(b *testing.B, l *zap.Logger) {
	b.Run("no fields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.Info("Simple log message")
		}
	})
	b.Run("application fields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.Info("Simple log message", zap.Uint32("count", 123), zap.String("details", "nothing"))
		}
	})
	b.Run("journald fields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.Info("Simple log message", SyslogIdentifier())
		}
	})
}
