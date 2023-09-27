package rollzap

import (
	"encoding/json"

	"github.com/pkg/errors"
	rollbar "github.com/rollbar/rollbar-go"
	"go.uber.org/zap/zapcore"
)

type RollbarCoreOptions struct {
	syncOnWrite bool
}

type RollbarCoreOption func(*RollbarCoreOptions)

func WithoutSyncOnWrite() RollbarCoreOption {
	return func(opt *RollbarCoreOptions) {
		opt.syncOnWrite = false
	}
}

// RollbarCore is a custom core to send logs to Rollbar. Add the core using zapcore.NewTee
type RollbarCore struct {
	zapcore.LevelEnabler

	coreFields  map[string]any
	syncOnWrite bool
}

// NewRollbarCore creates a new core to transmit logs to rollbar. rollbar token and other options should be set before creating a new core
func NewRollbarCore(minLevel zapcore.Level, options ...RollbarCoreOption) *RollbarCore {
	opts := &RollbarCoreOptions{
		syncOnWrite: true,
	}
	for _, opt := range options {
		opt(opts)
	}

	return &RollbarCore{
		LevelEnabler: minLevel,
		coreFields:   make(map[string]any),
		syncOnWrite:  opts.syncOnWrite,
	}
}

// With provides structure
func (c *RollbarCore) With(fields []zapcore.Field) zapcore.Core {
	fieldMap := fieldsToMap(fields)
	for k, v := range fieldMap {
		c.coreFields[k] = v
	}
	return c
}

// Check determines if this should be sent to roll bar based on LevelEnabler
func (c *RollbarCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

func (c *RollbarCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	fieldMap := fieldsToMap(fields)

	if len(c.coreFields) > 0 {
		coreFieldsMap, err := json.Marshal(c.coreFields)
		if err != nil {
			return errors.Wrapf(err, "unable to parse json for coreFields")
		}
		fieldMap["coreFields"] = string(coreFieldsMap)
	}

	if entry.LoggerName != "" {
		fieldMap["logger"] = entry.LoggerName
	}
	if entry.Caller.TrimmedPath() != "" {
		fieldMap["file"] = entry.Caller.TrimmedPath()
	}

	switch entry.Level {
	case zapcore.DebugLevel:
		rollbar.Debug(entry.Message, fieldMap)
	case zapcore.InfoLevel:
		rollbar.Info(entry.Message, fieldMap)
	case zapcore.WarnLevel:
		rollbar.Warning(entry.Message, fieldMap)
	case zapcore.ErrorLevel:
		errMap := extractError(fields)
		if errMap != nil {
			rollbar.Error(entry.Message, fieldMap, errMap)
		} else {
			rollbar.Error(entry.Message, fieldMap)
		}
	case zapcore.DPanicLevel:
		rollbar.Critical(entry.Message, fieldMap)
	case zapcore.PanicLevel:
		rollbar.Critical(entry.Message, fieldMap)
	case zapcore.FatalLevel:
		rollbar.Critical(entry.Message, fieldMap)

	}

	if c.syncOnWrite {
		rollbar.Wait()
	}

	return nil
}

// Sync flushes
func (c *RollbarCore) Sync() error {
	rollbar.Wait()
	return nil
}

func extractError(fields []zapcore.Field) error {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	var foundError error
	for _, f := range fields {
		if f.Type == zapcore.ErrorType {
			foundError = f.Interface.(error)
		}
	}
	return foundError
}

func fieldsToMap(fields []zapcore.Field) map[string]any {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	m := make(map[string]any)
	for k, v := range enc.Fields {
		m[k] = v
	}
	return m
}
