package logs

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LayeredError struct {
	Layers  []string
	Message string
}

type errorDetails struct {
	Message string
	Trace   string
	Err     error
}

func ParseErrorChain(err error) LayeredError {
	parts := strings.Split(err.Error(), " error: ")
	if len(parts) == 1 {
		return LayeredError{Message: err.Error()}
	}
	return LayeredError{
		Layers:  parts[:len(parts)-1],
		Message: parts[len(parts)-1],
	}
}

func parseError(message interface{}) errorDetails {
	switch v := message.(type) {
	case error:
		parsed := ParseErrorChain(v)
		msg := parsed.Message
		if msg == "" && len(parsed.Layers) > 0 {
			msg = parsed.Layers[0]
		}
		var trace string
		if len(parsed.Layers) > 0 {
			trace = strings.Join(parsed.Layers, " -> ")
		}
		return errorDetails{
			Message: msg,
			Trace:   trace,
			Err:     v,
		}
	case string:
		return errorDetails{Message: v}
	default:
		return errorDetails{Message: fmt.Sprintf("%v", v)}
	}
}

func printPrettyError(w io.Writer, pc uintptr, details errorDetails, args ...any) {
	ts := time.Now().Format(prettyTimeLayout)
	src := formatSource(pc)

	var b strings.Builder
	fmt.Fprintf(&b, "%s\tERROR\t%s\n", ts, src)
	fmt.Fprintf(&b, "Message: %s\n", details.Message)
	if details.Trace != "" {
		fmt.Fprintf(&b, "Trace: %s\n", details.Trace)
	}
	if len(args) > 0 {
		fmt.Fprintf(&b, "Fields: %s\n", formatArgs(args...))
	}
	_, _ = w.Write([]byte(b.String()))
}

func callerPC(skip int) uintptr {
	var pcs [1]uintptr
	runtime.Callers(skip, pcs[:])
	return pcs[0]
}

func formatSource(pc uintptr) string {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	if frame.File == "" {
		return "unknown"
	}
	file := filepath.ToSlash(frame.File)
	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		file = parts[len(parts)-2] + "/" + parts[len(parts)-1]
	} else {
		file = parts[len(parts)-1]
	}
	if frame.Line > 0 {
		return fmt.Sprintf("%s:%d", file, frame.Line)
	}
	return file
}

func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, (len(args)+1)/2)
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		if i+1 < len(args) {
			parts = append(parts, fmt.Sprintf("%s=%v", key, args[i+1]))
		} else {
			parts = append(parts, key)
		}
	}
	return strings.Join(parts, " ")
}
