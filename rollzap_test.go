package rollzap

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	rollbar "github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newEntry(level zapcore.Level, msg ...string) zapcore.Entry {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(fmt.Sprintf("Test %s Messsage", level))

	if len(msg) > 0 {
		strBuilder.WriteString(fmt.Sprintf(" %s", strings.Join(msg, " ")))
	}

	return zapcore.Entry{
		Level:      level,
		Time:       time.Now(),
		LoggerName: "test_logger",
		Message:    strBuilder.String(),
		Caller:     zapcore.EntryCaller{},
		Stack:      "No stack. Just test.",
	}
}

func TestRollbarCoreLevels(t *testing.T) {
	rc := NewRollbarCore(zapcore.DebugLevel)
	if rc == nil {
		t.Fatalf("RollbarCore not initialized")
	}
	if rc.LevelEnabler != zapcore.DebugLevel {
		t.Fatalf("rc.minLevel is not the correct value")
	}

	rc = NewRollbarCore(zapcore.ErrorLevel)
	if rc == nil {
		t.Fatalf("RollbarCore not initialized")
	}

	if rc.LevelEnabler != zapcore.ErrorLevel {
		t.Fatalf("rc.minLevel is not the correct value")
	}
}

func TestNewRollbarCore(t *testing.T) {
	token := os.Getenv("RC_TOKEN")
	if token == "" {
		t.Log("running tests without token")
	}
	rollbar.SetToken(token)

	rc := NewRollbarCore(zapcore.ErrorLevel)
	if rc == nil {
		t.Fatalf("RollbarCore not initialized")
	}

	if rc.LevelEnabler != zapcore.ErrorLevel {
		t.Fatalf("rc.minLevel is not the correct value")
	}

	coreFields := []zap.Field{zap.String("foo", "bar"), zap.String("moo", "cow")}
	rc.With(coreFields)
	if _, ok := rc.coreFields["foo"]; !ok {
		t.Fatalf("core fields not stored")
	}

	if _, ok := rc.coreFields["moo"]; !ok {
		t.Fatalf("core fields not stored")
	}

	debugEntry := newEntry(zapcore.DebugLevel)
	rc.Check(debugEntry, nil) // should do more here?
	if err := rc.Write(debugEntry, nil); err != nil {
		t.Errorf("writing debug message %v", err)
	}

	// all of these will write, there are no checks, but in a real scenario, they work. Probably need to revisit this test
	// Right now, this is just to make sure the different levels are entered into Rollbar correctly.
	if err := rc.Write(newEntry(zapcore.InfoLevel), nil); err != nil {
		t.Errorf("writing Info message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.WarnLevel), nil); err != nil {
		t.Errorf("writing warn message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.ErrorLevel), nil); err != nil {
		t.Errorf("writing error message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.DPanicLevel), nil); err != nil {
		t.Errorf("writing dpanic message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.PanicLevel), nil); err != nil {
		t.Errorf("writing panic message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.FatalLevel), nil); err != nil {
		t.Errorf("writing fatal message %v", err)
	}

	if err := rc.Sync(); err != nil {
		t.Errorf("sync failed - %v", err)
	}
}

func TestNewRollbarCoreWithoutSyncOnWrite(t *testing.T) {
	token := os.Getenv("RC_TOKEN")
	if token == "" {
		t.Log("running tests without token")
	}
	rollbar.SetToken(token)

	rc := NewRollbarCore(zapcore.ErrorLevel, WithoutSyncOnWrite())
	if rc == nil {
		t.Fatalf("RollbarCore not initialized")
	}

	if rc.LevelEnabler != zapcore.ErrorLevel {
		t.Fatalf("rc.minLevel is not the correct value")
	}

	coreFields := []zap.Field{zap.String("foo", "bar"), zap.String("moo", "cow")}
	rc.With(coreFields)
	if _, ok := rc.coreFields["foo"]; !ok {
		t.Fatalf("core fields not stored")
	}

	if _, ok := rc.coreFields["moo"]; !ok {
		t.Fatalf("core fields not stored")
	}

	debugEntry := newEntry(zapcore.DebugLevel, "without SyncOnWrite")
	rc.Check(debugEntry, nil) // should do more here?
	if err := rc.Write(debugEntry, nil); err != nil {
		t.Errorf("writing debug message %v", err)
	}

	// all of these will write, there are no checks, but in a real scenario, they work. Probably need to revisit this test
	// Right now, this is just to make sure the different levels are entered into Rollbar correctly.
	if err := rc.Write(newEntry(zapcore.InfoLevel, "without SyncOnWrite"), nil); err != nil {
		t.Errorf("writing Info message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.WarnLevel, "without SyncOnWrite"), nil); err != nil {
		t.Errorf("writing warn message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.ErrorLevel, "without SyncOnWrite"), nil); err != nil {
		t.Errorf("writing error message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.DPanicLevel, "without SyncOnWrite"), nil); err != nil {
		t.Errorf("writing dpanic message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.PanicLevel, "without SyncOnWrite"), nil); err != nil {
		t.Errorf("writing panic message %v", err)
	}
	if err := rc.Write(newEntry(zapcore.FatalLevel, "without SyncOnWrite"), nil); err != nil {
		t.Errorf("writing fatal message %v", err)
	}

	if err := rc.Sync(); err != nil {
		t.Errorf("sync failed - %v", err)
	}
}
