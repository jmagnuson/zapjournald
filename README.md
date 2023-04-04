# Zap Core for systemd Journal
Package `zapjournld` provides `zap` Core for systemd Journal. It supports structured logging.

Applications may use zap logger to write logs directly into journald and may relatively freely define additional fields, which will be indexed by journald.

Information about common and custom journald fields is available on https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html.
## Example
```go
package main

import (
	"fmt"

	"git.frostfs.info/TrueCloudLab/zapjournald"
	"go.uber.org/zap"

	"github.com/ssgreg/journald"
	"go.uber.org/zap/zapcore"
)

func main() {
	// StandardLogger
	standardLogger, _ := NewStandardLogger("info")

	standardLogger.Info("Simple log raw 1")
}

func NewStandardLogger(lvlStr string) (*zap.Logger, zap.AtomicLevel) {
	lvl, err := getLogLevel(lvlStr)
	if err != nil {
		panic(err)
	}

	zc := zap.NewProductionConfig()
	zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zc.Level = zap.NewAtomicLevelAt(lvl)

	// Initialize Zap.
	encoder := zapcore.NewJSONEncoder(zc.EncoderConfig)

	core := zapjournald.NewCore(zap.NewAtomicLevelAt(lvl), encoder, &journald.Journal{}, zapjournald.SyslogFields)
	coreWithContext := core.With([]zapcore.Field{
		zapjournald.SyslogFacility(zapjournald.LogDaemon),
		zapjournald.SyslogIdentifier(),
		zapjournald.SyslogPid(),
	})
	l := zap.New(coreWithContext, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.FatalLevel)))
	return l, zc.Level
}

func getLogLevel(lvlStr string) (zapcore.Level, error) {
	var lvl zapcore.Level
	err := lvl.UnmarshalText([]byte(lvlStr))
	if err != nil {
		return lvl, fmt.Errorf("incorrect logger level configuration %s (%v), "+
			"value should be one of %v", lvlStr, err, [...]zapcore.Level{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			zapcore.WarnLevel,
			zapcore.ErrorLevel,
			zapcore.DPanicLevel,
			zapcore.PanicLevel,
			zapcore.FatalLevel,
		})
	}
	return lvl, nil
}
```

### Results

#### You can read simple view like this:
```
sudo journalctl -f
Apr 07 13:06:34 user ___11go_build_git_frostfs_info_TrueCloudLab_zapjournald_main[330983]: {"level":"info","ts":"2023-04-07T13:06:34.990+0300","msg":"Simple log raw 1","SYSLOG_FACILITY":3,"SYSLOG_IDENTIFIER":"___11go_build_git_frostfs_info_TrueCloudLab_zapjournald_main","SYSLOG_PID":330983}
```
#### Or you can find lines by indexed field
```
sudo journalctl SYSLOG_PID=76385
Apr 07 13:06:34 user ___11go_build_git_frostfs_info_TrueCloudLab_zapjournald_main[330983]: {"level":"info","ts":"2023-04-07T13:06:34.990+0300","msg":"Simple log raw 1","SYSLOG_FACILITY":3,"SYSLOG_ID>
```
#### Or you can read full-structured view like this:
```
sudo journalctl SYSLOG_PID=76385 --output=json-pretty
{
        "PRIORITY" : "6",
        "_SYSTEMD_OWNER_UID" : "1000",
        "_HOSTNAME" : "user",
        "__MONOTONIC_TIMESTAMP" : "153798818198",
        "_SYSTEMD_SLICE" : "user-1000.slice",
        "_GID" : "1000",
        "_AUDIT_LOGINUID" : "1000",
        "_UID" : "1000",
        "SYSLOG_IDENTIFIER" : "___11go_build_git_frostfs_info_TrueCloudLab_zapjournald_main",
        "__REALTIME_TIMESTAMP" : "1680861994990646",
        "_CAP_EFFECTIVE" : "0",
        "_COMM" : "___11go_build_g",
        "_AUDIT_SESSION" : "59",
        "__CURSOR" : "s=9e2d157a286a437ea3405618239c1e07;i=14a96;b=bc1bc632f462476abc9e3e3b0c517f6c;m=23cf1fb996;t=5f8bc2e20bc36;x=9c6c4a86b6da43c",
        "_SYSTEMD_USER_SLICE" : "app.slice",
        "_SYSTEMD_CGROUP" : "/user.slice/user-1000.slice/user@1000.service/app.slice/app-gnome-jetbrains\\x2dgoland-203587.scope",
        "MESSAGE" : "{\"level\":\"info\",\"ts\":\"2023-04-07T13:06:34.990+0300\",\"msg\":\"Simple log raw 1\",\"SYSLOG_FACILITY\":3,\"SYSLOG_IDENTIFIER\":\"___11go_build_git_frostfs_info_TrueCloudLab_zapjournal>
        "_SOURCE_REALTIME_TIMESTAMP" : "1680861994990539",
        "_MACHINE_ID" : "82be99c92deb4b36850366d7742de698",
        "SYSLOG_FACILITY" : "3",
        "_BOOT_ID" : "bc1bc632f462476abc9e3e3b0c517f6c",
        "_PID" : "330983",
        "_SELINUX_CONTEXT" : "unconfined\n",
        "_SYSTEMD_UNIT" : "user@1000.service",
        "_TRANSPORT" : "journal",
        "_SYSTEMD_INVOCATION_ID" : "ed23dc57a96340029ef8bf9956182c20",
        "_SYSTEMD_USER_UNIT" : "app-gnome-jetbrains\\x2dgoland-203587.scope",
        "SYSLOG_PID" : "330983"
}
```
