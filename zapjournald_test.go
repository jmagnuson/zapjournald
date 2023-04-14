package zapjournald

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestGetFieldValueWithStringFormatting testing unexported function getFieldValue and verifying that all zap field types have correct return values.
func TestGetFieldValueWithStringFormatting(t *testing.T) {
	// Interface types.
	addr := net.ParseIP("1.2.3.4")
	name := usernameTestExample("phil")
	ints := []int{5, 6}

	tests := []struct {
		name   string
		field  zap.Field
		expect interface{}
	}{
		{"ObjectMarshaler", zap.Object("k", name), "phil"},
		{"ArrayMarshaler", zap.Array("k", boolsTestExample([]bool{true})), "[true]"},
		{"Binary", zap.Binary("k", []byte("ab12")), "ab12"},
		{"Bool", zap.Bool("k", true), "1"},
		{"ByteString", zap.ByteString("k", []byte("ab12")), "ab12"},
		{"Complex128", zap.Complex128("k", 1+2i), "(1+2i)"},
		{"Complex64", zap.Complex64("k", 1+2i), "(1+2i)"},
		{"Duration", zap.Duration("k", 1), "1"},
		{"Float64", zap.Float64("k", 3.14), "4614253070214989087"},
		{"Float32", zap.Float32("k", 3.14), "1078523331"},
		{"Int64", zap.Int64("k", 1), "1"},
		{"Int32", zap.Int32("k", 1), "1"},
		{"Int16", zap.Int16("k", 1), "1"},
		{"Int8", zap.Int8("k", 1), "1"},
		{"String", zap.String("k", "foo"), "foo"},
		{"Time", zap.Time("k", time.Unix(0, 0).In(time.UTC)), "0 UTC"},
		{"TimeFull", zap.Time("k", time.Time{}), "0001-01-01 00:00:00 +0000 UTC"},
		{"Uint", zap.Uint("k", 1), "1"},
		{"Uint64", zap.Uint64("k", 1), "1"},
		{"Uint32", zap.Uint32("k", 1), "1"},
		{"Uint16", zap.Uint16("k", 1), "1"},
		{"Uint8", zap.Uint8("k", 1), "1"},
		{"Uintptr", zap.Uintptr("k", 0xa), "10"},
		{"Namespace", zap.Namespace("k"), "<nil>"},
		{"Stringer", zap.Stringer("k", addr), "1.2.3.4"},
		{"Reflect", zap.Reflect("k", ints), "[5 6]"},
		{"Error", zap.Error(fmt.Errorf("test")), "test"},
		{"Skip", zap.Skip(), "<nil>"},
		{"InlineMarshaller", zap.Inline(name), "phil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldValue := valueToString(getFieldValue(tt.field))
			if !assert.Equal(t, tt.expect, fieldValue, "Unexpected output  %s.", tt.name) {
				t.Logf("type expected: %T\nGot: %T", tt.expect, fieldValue)
			}
		})
	}
}

// usernameTestExample type, which are implements ObjectMarshaller interface.
type usernameTestExample string

func (n usernameTestExample) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}

// boolsTestExample type, which are implements ArrayMarshaller interface.
type boolsTestExample []bool

func (bs boolsTestExample) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bs {
		arr.AppendBool(bs[i])
	}
	return nil
}

func valueToString(value interface{}) string {
	switch rv := value.(type) {
	case string:
		return rv
	case []byte:
		return string(rv)
	default:
		return fmt.Sprint(value)
	}
}
