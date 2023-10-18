package zapjournald

import (
	"bytes"
	"fmt"
	"strings"

	"go.uber.org/zap/buffer"
)

var pool = buffer.NewPool()

func encodeJournaldField(buf *buffer.Buffer, key string, value any) {
	switch v := value.(type) {
	case string:
		writeField(buf, key, v)
	case []byte:
		writeFieldBytes(buf, key, v)
	default:
		writeField(buf, key, fmt.Sprint(v))
	}
}

func writeFieldBytes(buf *buffer.Buffer, name string, value []byte) {
	buf.Write([]byte(name))
	if bytes.ContainsRune(value, '\n') {
		// According to the format, if the value includes a newline
		// need to write the field name, plus a newline, then the
		// size (64bit LE), the field data and a final newline.

		buf.Write([]byte{'\n'})
		appendUint64Binary(buf, uint64(len(value)))
	} else {
		buf.Write([]byte{'='})
	}
	buf.Write(value)
	buf.Write([]byte{'\n'})
}

func writeField(buf *buffer.Buffer, name string, value string) {
	buf.Write([]byte(name))
	if strings.ContainsRune(value, '\n') {
		// According to the format, if the value includes a newline
		// need to write the field name, plus a newline, then the
		// size (64bit LE), the field data and a final newline.

		buf.Write([]byte{'\n'})
		// 1 allocation here.
		// binary.Write(w, binary.LittleEndian, uint64(len(value)))
		appendUint64Binary(buf, uint64(len(value)))
	} else {
		buf.Write([]byte{'='})
	}
	buf.WriteString(value)
	buf.Write([]byte{'\n'})
}

func appendUint64Binary(buf *buffer.Buffer, v uint64) {
	// Copied from https://github.com/golang/go/blob/go1.21.3/src/encoding/binary/binary.go#L119
	buf.Write([]byte{
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
		byte(v >> 24),
		byte(v >> 32),
		byte(v >> 40),
		byte(v >> 48),
		byte(v >> 56),
	})
}
