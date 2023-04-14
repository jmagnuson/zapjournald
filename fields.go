package zapjournald

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Facility defines syslog facility from /usr/include/sys/syslog.h .
type Facility uint32

// Syslog compatibility fields.
const (
	SyslogFacilityField   = "SYSLOG_FACILITY"
	SyslogIdentifierField = "SYSLOG_IDENTIFIER"
	SyslogPidField        = "SYSLOG_PID"
	SyslogTimestampField  = "SYSLOG_TIMESTAMP"
)

// Facilities before LogFtp are the same on Linux, BSD, and OS X.
const (
	LogKern Facility = iota << 3
	LogUser
	LogMail
	LogDaemon
	LogAuth
	LogSyslog
	LogLpr
	LogNews
	LogUucp
	LogCron
	LogAuthpriv
	LogFtp
	_        // unused
	LogAudit // unused
	_        // unused
	_        // unused
	LogLocal0
	LogLocal1
	LogLocal2
	LogLocal3
	LogLocal4
	LogLocal5
	LogLocal6
	LogLocal7
)

// LogFacmask is used to extract facility part of the message.
const LogFacmask = 0x03f8

// SyslogFields contains slice of fields that are
// indexed by syslog by default.
var SyslogFields = []string{
	SyslogFacilityField,
	SyslogIdentifierField,
	SyslogPidField,
	SyslogTimestampField,
}

func SyslogFacility(facility Facility) zapcore.Field {
	return zap.Uint32(SyslogFacilityField, uint32(facility&LogFacmask)>>3)
}

func SyslogIdentifier() zapcore.Field {
	return zap.String(SyslogIdentifierField, filepath.Base(os.Args[0]))
}

func SyslogPid() zapcore.Field {
	return zap.Int(SyslogPidField, os.Getpid())
}

func SyslogTimestamp(time time.Time) zapcore.Field {
	return zap.Time(SyslogTimestampField, time)
}
