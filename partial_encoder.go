package zapjournald

import (
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"
)

type partialEncoder struct {
	wrap   zapcore.Encoder
	ignore map[string]struct{}
}

// NewPartialEncoder wraps existing encoder to avoid output of some provided
// fields. The main use case is to ignore SyslogFields that leak into
// ConsoleEncoder and provide no additional info for the human.
func NewPartialEncoder(enc zapcore.Encoder, ignore []string) zapcore.Encoder {
	m := make(map[string]struct{}, len(ignore))
	for _, i := range ignore {
		m[i] = struct{}{}
	}
	return partialEncoder{
		wrap:   enc,
		ignore: m,
	}
}

func (enc partialEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	if _, ok := enc.ignore[key]; ok {
		return nil
	}
	return enc.wrap.AddArray(key, marshaler)
}

func (enc partialEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	if _, ok := enc.ignore[key]; ok {
		return nil
	}
	return enc.wrap.AddObject(key, marshaler)
}

func (enc partialEncoder) AddBinary(key string, value []byte) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddBinary(key, value)
}

func (enc partialEncoder) AddByteString(key string, value []byte) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddByteString(key, value)
}

func (enc partialEncoder) AddBool(key string, value bool) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddBool(key, value)
}

func (enc partialEncoder) AddComplex128(key string, value complex128) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddComplex128(key, value)
}

func (enc partialEncoder) AddComplex64(key string, value complex64) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddComplex64(key, value)
}

func (enc partialEncoder) AddDuration(key string, value time.Duration) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddDuration(key, value)
}

func (enc partialEncoder) AddFloat64(key string, value float64) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddFloat64(key, value)
}

func (enc partialEncoder) AddFloat32(key string, value float32) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddFloat32(key, value)
}

func (enc partialEncoder) AddInt(key string, value int) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddInt(key, value)
}

func (enc partialEncoder) AddInt64(key string, value int64) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddInt64(key, value)
}

func (enc partialEncoder) AddInt32(key string, value int32) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddInt32(key, value)
}

func (enc partialEncoder) AddInt16(key string, value int16) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddInt16(key, value)
}

func (enc partialEncoder) AddInt8(key string, value int8) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddInt8(key, value)
}

func (enc partialEncoder) AddString(key, value string) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddString(key, value)
}

func (enc partialEncoder) AddTime(key string, value time.Time) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddTime(key, value)
}

func (enc partialEncoder) AddUint(key string, value uint) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddUint(key, value)
}

func (enc partialEncoder) AddUint64(key string, value uint64) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddUint64(key, value)
}

func (enc partialEncoder) AddUint32(key string, value uint32) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddUint32(key, value)
}

func (enc partialEncoder) AddUint16(key string, value uint16) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddUint16(key, value)
}

func (enc partialEncoder) AddUint8(key string, value uint8) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddUint8(key, value)
}

func (enc partialEncoder) AddUintptr(key string, value uintptr) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.AddUintptr(key, value)
}

func (enc partialEncoder) AddReflected(key string, value interface{}) error {
	if _, ok := enc.ignore[key]; ok {
		return nil
	}
	return enc.wrap.AddReflected(key, value)
}

func (enc partialEncoder) OpenNamespace(key string) {
	if _, ok := enc.ignore[key]; ok {
		return
	}
	enc.wrap.OpenNamespace(key)
}

func (enc partialEncoder) Clone() zapcore.Encoder {
	return partialEncoder{
		wrap:   enc.wrap.Clone(),
		ignore: maps.Clone(enc.ignore),
	}
}

func (enc partialEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return enc.wrap.EncodeEntry(entry, fields)
}
