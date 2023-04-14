package zapjournald

import (
	"fmt"

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
	// Generate the message.
	buffer, err := core.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return fmt.Errorf("failed to encode log entry: %w", err)
	}

	message := buffer.String()

	structuredFields := maps.Clone(core.contextStructuredFields)
	for _, field := range fields {
		if _, isJournalField := core.storedFieldNames[field.Key]; isJournalField {
			structuredFields[field.Key] = getFieldValue(field)
		}
	}

	// Write the message.
	switch entry.Level {
	case zapcore.DebugLevel:
		return core.j.Send(message, journald.PriorityDebug, structuredFields)

	case zapcore.InfoLevel:
		return core.j.Send(message, journald.PriorityInfo, structuredFields)

	case zapcore.WarnLevel:
		return core.j.Send(message, journald.PriorityWarning, structuredFields)

	case zapcore.ErrorLevel:
		return core.j.Send(message, journald.PriorityErr, structuredFields)

	case zapcore.DPanicLevel:
		return core.j.Send(message, journald.PriorityCrit, structuredFields)

	case zapcore.PanicLevel:
		return core.j.Send(message, journald.PriorityCrit, structuredFields)

	case zapcore.FatalLevel:
		return core.j.Send(message, journald.PriorityCrit, structuredFields)

	default:
		return fmt.Errorf("unknown log level: %v", entry.Level)
	}
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
	case zapcore.DurationType,
		zapcore.Float64Type,
		zapcore.Float32Type,
		zapcore.Int64Type,
		zapcore.Int32Type,
		zapcore.Int16Type,
		zapcore.Int8Type,
		zapcore.Uint64Type,
		zapcore.Uint32Type,
		zapcore.Uint16Type,
		zapcore.Uint8Type,
		zapcore.UintptrType,
		zapcore.BoolType:
		return f.Integer
	case zapcore.StringType:
		return f.String
	case zapcore.TimeType:
		if f.Interface != nil {
			// for example: zap.Time("k", time.Unix(100900, 0).In(time.UTC)) - will produce: "100900000000000 UTC" (result in nanoseconds)
			return fmt.Sprintf("%d %v", f.Integer, f.Interface)
		}
		return f.Integer
	default:
		panic(fmt.Sprintf("unknown field type: %v", f))
	}
}
