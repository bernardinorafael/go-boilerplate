package logging

import (
	"time"

	"go.uber.org/zap"
)

// LogField is an interface that represents a generic log field
type LogField interface {
	ToZapField() zap.Field
}

// Helper functions to create log fields

func String(key, value string) LogField {
	return StringField{Key: key, Value: value}
}

func Int(key string, value int) LogField {
	return IntField{Key: key, Value: value}
}

func Int32(key string, value int32) LogField {
	return Int32Field{Key: key, Value: value}
}

func Int64(key string, value int64) LogField {
	return Int64Field{Key: key, Value: value}
}

func Uint(key string, value uint) LogField {
	return UintField{Key: key, Value: value}
}

func Uint32(key string, value uint32) LogField {
	return Uint32Field{Key: key, Value: value}
}

func Uint64(key string, value uint64) LogField {
	return Uint64Field{Key: key, Value: value}
}

func Float32(key string, value float32) LogField {
	return Float32Field{Key: key, Value: value}
}

func Float64(key string, value float64) LogField {
	return Float64Field{Key: key, Value: value}
}

func Bool(key string, value bool) LogField {
	return BoolField{Key: key, Value: value}
}

func Time(key string, value time.Time) LogField {
	return TimeField{Key: key, Value: value}
}

func Duration(key string, value time.Duration) LogField {
	return DurationField{Key: key, Value: value}
}

func Err(err error) LogField {
	return ErrorField{Err: err}
}

func Any(key string, value interface{}) LogField {
	return AnyField{Key: key, Value: value}
}

// Binary creates a field that carries an opaque binary blob.
func Binary(key string, value []byte) LogField {
	return AnyField{Key: key, Value: value}
}

// StringField represents a log field of type string
type StringField struct {
	Key   string
	Value string
}

func (f StringField) ToZapField() zap.Field {
	return zap.String(f.Key, f.Value)
}

// IntField represents a log field of type int
type IntField struct {
	Key   string
	Value int
}

func (f IntField) ToZapField() zap.Field {
	return zap.Int(f.Key, f.Value)
}

// Int32Field represents a log field of type int32
type Int32Field struct {
	Key   string
	Value int32
}

func (f Int32Field) ToZapField() zap.Field {
	return zap.Int32(f.Key, f.Value)
}

// Int64Field represents a log field of type int64
type Int64Field struct {
	Key   string
	Value int64
}

func (f Int64Field) ToZapField() zap.Field {
	return zap.Int64(f.Key, f.Value)
}

// UintField represents a log field of type uint
type UintField struct {
	Key   string
	Value uint
}

func (f UintField) ToZapField() zap.Field {
	return zap.Uint(f.Key, f.Value)
}

// Uint32Field represents a log field of type uint32
type Uint32Field struct {
	Key   string
	Value uint32
}

func (f Uint32Field) ToZapField() zap.Field {
	return zap.Uint32(f.Key, f.Value)
}

// Uint64Field represents a log field of type uint64
type Uint64Field struct {
	Key   string
	Value uint64
}

func (f Uint64Field) ToZapField() zap.Field {
	return zap.Uint64(f.Key, f.Value)
}

// Float32Field represents a log field of type float32
type Float32Field struct {
	Key   string
	Value float32
}

func (f Float32Field) ToZapField() zap.Field {
	return zap.Float32(f.Key, f.Value)
}

// Float64Field represents a log field of type float64
type Float64Field struct {
	Key   string
	Value float64
}

func (f Float64Field) ToZapField() zap.Field {
	return zap.Float64(f.Key, f.Value)
}

// BoolField represents a log field of type bool
type BoolField struct {
	Key   string
	Value bool
}

func (f BoolField) ToZapField() zap.Field {
	return zap.Bool(f.Key, f.Value)
}

// TimeField represents a log field of type time.Time
type TimeField struct {
	Key   string
	Value time.Time
}

func (f TimeField) ToZapField() zap.Field {
	return zap.Time(f.Key, f.Value)
}

// DurationField represents a log field of type time.Duration
type DurationField struct {
	Key   string
	Value time.Duration
}

func (f DurationField) ToZapField() zap.Field {
	return zap.Duration(f.Key, f.Value)
}

// ErrorField represents a log field of type error
type ErrorField struct {
	Err error
}

func (f ErrorField) ToZapField() zap.Field {
	return zap.Error(f.Err)
}

// AnyField represents a log field of any type
type AnyField struct {
	Key   string
	Value interface{}
}

func (f AnyField) ToZapField() zap.Field {
	return zap.Any(f.Key, f.Value)
}
