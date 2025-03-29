package logging

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type customJSONFormatter struct {
	logToFile bool
	attr      []zap.Field
}

func newCustomJSONFormatter(params LogParams) *customJSONFormatter {
	res := &customJSONFormatter{
		logToFile: params.LogToFile,
	}

	if params.AppName != "" {
		res.attr = append(res.attr, zap.String("app", params.AppName))
	}

	return res
}

func (f *customJSONFormatter) formatLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	levelStr := l.CapitalString()
	if l == LevelCritical {
		levelStr = "CRITICAL"
	}
	enc.AppendString(levelStr)
}

var levelColors = map[string]func(...any) string{
	"DEBUG":    color.Magenta.Render,
	"INFO":     color.Blue.Render,
	"WARN":     color.Yellow.Render,
	"ERROR":    color.Red.Render,
	"FATAL":    func(a ...any) string { return color.Bold.Render(color.Red.Render(a...)) },
	"CRITICAL": func(a ...any) string { return color.Bold.Render(color.Red.Render(a...)) },
}

func (f *customJSONFormatter) getRuntimeData() (funcName, filename string, line int) {
	pc, filePath, line, ok := runtime.Caller(9)
	if !ok {
		return "unknown", "unknown", 0
	}
	filename = filepath.Base(filePath)
	funcPath := runtime.FuncForPC(pc).Name()
	funcName = funcPath[strings.LastIndex(funcPath, ".")+1:]

	if strings.Contains(funcName, "func") {
		funcBefore := funcPath[:strings.LastIndex(funcPath, ".")]
		funcName = funcPath[strings.LastIndex(funcBefore, ".")+1:]
	}
	return
}

func (f *customJSONFormatter) getDefaultFields() []zap.Field {
	return f.attr
}

func (f *customJSONFormatter) Format(logEntry string) string {
	if f.logToFile {
		return logEntry
	}

	var entry map[string]any
	if err := json.Unmarshal([]byte(logEntry), &entry); err != nil {
		return logEntry
	}

	entry = f.applyColorToLevel(entry)
	entry = f.formatRuntimeDetails(entry)
	return f.reorderFields(entry)
}

func (f *customJSONFormatter) applyColorToLevel(entry map[string]any) map[string]any {
	level, ok := entry["level"].(string)
	if ok {
		if colorFunc, exists := levelColors[level]; exists {
			entry["level"] = colorFunc(level)
		}
	}
	return entry
}

func (f *customJSONFormatter) formatRuntimeDetails(entry map[string]any) map[string]any {
	if msg, ok := entry["msg"].(string); ok {
		funcName, _, _ := f.getRuntimeData()
		entry["msg"] = fmt.Sprintf("%s: %s", funcName, msg)
	}

	return entry
}

func (f *customJSONFormatter) reorderFields(entry map[string]any) string {
	orderedFields := []string{"time", "level", "msg", "file", "app", "host"}
	var result strings.Builder

	result.Grow(len(entry) * 10)

	result.WriteString("{")
	firstField := true

	for _, field := range orderedFields {
		if value, exists := entry[field]; exists {
			if !firstField {
				result.WriteString(",")
			}
			firstField = false
			writeField(&result, field, value, true)
			delete(entry, field)
		}
	}

	for key, value := range entry {
		if !firstField {
			result.WriteString(",")
		}
		firstField = false
		writeField(&result, key, value, false)
	}

	result.WriteString("}\n")
	return result.String()
}

func writeField(b *strings.Builder, key string, value any, isKnownString bool) {
	b.WriteString(`"`)
	b.WriteString(key)
	b.WriteString(`":`)

	if isKnownString {
		b.WriteString(`"`)
		b.WriteString(fmt.Sprint(value))
		b.WriteString(`"`)
		return
	}

	// common primitive values do not need to do json.Marshal
	switch v := value.(type) {
	case string:
		b.WriteString(`"`)
		b.WriteString(v)
		b.WriteString(`"`)
	case int, int8, int16, int32, int64:
		b.WriteString(strconv.FormatInt(v.(int64), 10))
	case uint, uint8, uint16, uint32, uint64:
		b.WriteString(strconv.FormatUint(v.(uint64), 10))
	case float32, float64:
		b.WriteString(strconv.FormatFloat(v.(float64), 'f', -1, 64))
	case bool:
		b.WriteString(strconv.FormatBool(v))
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			b.WriteString(`"`)
			b.WriteString(fmt.Sprint(value))
			b.WriteString(`"`)
		} else {
			b.Write(jsonValue)
		}
	}
}
