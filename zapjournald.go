package zapjournald

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ssgreg/journald"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"
)

// Core for zapjournald.
//
// Implements zapcore.LevelEnabler and zapcore.Core interfaces.
type Core struct {
	zapcore.LevelEnabler
	encoder zapcore.Encoder
	j       *journald.Journal
	// field names, which will be stored in journald structure
	storedFieldNames map[string]struct{}
	// journald fields, which are always present in current core context
	contextStructuredFields map[string]interface{}
}

func NewCore(enab zapcore.LevelEnabler, encoder zapcore.Encoder, journal *journald.Journal, journalFields []string) *Core {
	journalFieldsMap := make(map[string]struct{})
	for _, field := range journalFields {
		journalFieldsMap[field] = struct{}{}
	}
	return &Core{
		LevelEnabler:            enab,
		encoder:                 encoder,
		j:                       journal,
		storedFieldNames:        journalFieldsMap,
		contextStructuredFields: make(map[string]interface{}),
	}
}

// With adds structured context to the Core.
func (core *Core) With(fields []zapcore.Field) zapcore.Core {
	clone := core.clone()
	for _, field := range fields {
		field.AddTo(clone.encoder)
		clone.contextStructuredFields[field.Key] = getFieldValue(field)
	}

	return clone
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
//
// Callers must use Check before calling Write.
func (core *Core) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if core.Enabled(entry.Level) {
		return checked.AddCore(entry, core)
	}
	return checked
}

// Write serializes the Entry and any Fields supplied at the log site and
// writes them to their destination.
//
// If called, Write should always log the Entry and Fields; it should not
// replicate the logic of Check.
func (core *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	prio, err := zapLevelToJournald(entry.Level)
	if err != nil {
		return err
	}

	b := pool.Get()
	defer b.Free()

	writeField(b, "PRIORITY", strconv.Itoa(int(prio)))

	if len(core.contextStructuredFields) != 0 {
		for k, v := range core.contextStructuredFields {
			encodeJournaldField(b, k, v)
		}
		for _, field := range fields {
			if _, isJournalField := core.storedFieldNames[field.Key]; isJournalField {
				encodeJournaldField(b, field.Key, getFieldValue(field))
			}
		}
	}

	// Generate the message.
	buffer, err := core.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return fmt.Errorf("failed to encode log entry: %w", err)
	}
	defer buffer.Free()

	writeFieldBytes(b, "MESSAGE", buffer.Bytes())

	// Write the message.
	return core.j.WriteMsg(b.Bytes())
}

// Sync flushes buffered logs (not used).
func (core *Core) Sync() error {
	return nil
}

// clone returns clone of core.
func (core *Core) clone() *Core {
	return &Core{
		LevelEnabler:            core.LevelEnabler,
		encoder:                 core.encoder.Clone(),
		j:                       core.j,
		storedFieldNames:        maps.Clone(core.storedFieldNames),
		contextStructuredFields: maps.Clone(core.contextStructuredFields),
	}
}

// getFieldValue returns underlying value stored in zapcore.Field.
func getFieldValue(f zapcore.Field) interface{} {
	switch f.Type {
	case zapcore.ArrayMarshalerType,
		zapcore.ObjectMarshalerType,
		zapcore.InlineMarshalerType,
		zapcore.BinaryType,
		zapcore.ByteStringType,
		zapcore.Complex128Type,
		zapcore.Complex64Type,
		zapcore.TimeFullType,
		zapcore.ReflectType,
		zapcore.NamespaceType,
		zapcore.StringerType,
		zapcore.ErrorType,
		zapcore.SkipType:
		return f.Interface
	case zapcore.DurationType:
		return time.Duration(f.Integer).String()
	case zapcore.Float64Type:
		// See https://github.com/uber-go/zap/blob/v1.26.0/buffer/buffer.go#L79
		f := math.Float64frombits(uint64(f.Integer))
		return strconv.FormatFloat(f, 'f', -1, 64)
	case zapcore.Float32Type:
		f := math.Float32frombits(uint32(f.Integer))
		return strconv.FormatFloat(float64(f), 'f', -1, 32)
	case zapcore.Int64Type,
		zapcore.Int32Type,
		zapcore.Int16Type,
		zapcore.Int8Type:
		return strconv.FormatInt(f.Integer, 10)
	case
		zapcore.Uint64Type,
		zapcore.Uint32Type,
		zapcore.Uint16Type,
		zapcore.Uint8Type,
		zapcore.UintptrType:
		return strconv.FormatUint(uint64(f.Integer), 10)
	case zapcore.BoolType:
		return strconv.FormatBool(f.Integer == 1)
	case zapcore.StringType:
		return f.String
	case zapcore.TimeType:
		if f.Interface != nil {
			// for example: zap.Time("k", time.Unix(100900, 0).In(time.UTC)) - will produce: "100900000000000 UTC" (result in nanoseconds)
			return fmt.Sprintf("%d %v", f.Integer, f.Interface)
		}
		return strconv.FormatUint(uint64(f.Integer), 10)
	default:
		panic(fmt.Sprintf("unknown field type: %v", f))
	}
}

func zapLevelToJournald(l zapcore.Level) (journald.Priority, error) {
	switch l {
	case zapcore.DebugLevel:
		return journald.PriorityDebug, nil
	case zapcore.InfoLevel:
		return journald.PriorityInfo, nil
	case zapcore.WarnLevel:
		return journald.PriorityWarning, nil
	case zapcore.ErrorLevel:
		return journald.PriorityErr, nil
	case zapcore.DPanicLevel:
		return journald.PriorityCrit, nil
	case zapcore.PanicLevel:
		return journald.PriorityCrit, nil
	case zapcore.FatalLevel:
		return journald.PriorityCrit, nil
	default:
		return 0, fmt.Errorf("unknown log level: %v", l)
	}
}
