package logs

import (
	"fmt"
	"io"
	"os"
	"time"

	dump "gofiber-hax/pkg/dump"
)

func printPrettyDebug(w io.Writer, pc uintptr, details errorDetails, message interface{}, args ...any) {
	ts := time.Now().Format(prettyTimeLayout)
	src := formatSource(pc)
	msg := details.Message
	if !isSimpleMessage(message) {
		msg = ""
	}

	if msg != "" {
		fmt.Fprintf(w, "%s\tDEBUG\t%s\t%s\n", ts, src, msg)
	} else {
		fmt.Fprintf(w, "%s\tDEBUG\t%s\n", ts, src)
	}

	payload := debugPayload(message, args...)
	if payload != nil {
		dump.Print(payload)
	}
}

func debugPayload(message interface{}, args ...any) any {
	if !isSimpleMessage(message) {
		return message
	}
	if len(args) == 0 {
		return nil
	}
	if len(args) == 1 {
		return args[0]
	}
	return kvArgsToMap(args...)
}

func kvArgsToMap(args ...any) map[string]any {
	if len(args) == 0 {
		return nil
	}
	m := make(map[string]any, (len(args)+1)/2)
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		if i+1 < len(args) {
			m[key] = args[i+1]
		} else {
			m[key] = true
		}
	}
	return m
}

func isSimpleMessage(message interface{}) bool {
	switch message.(type) {
	case string, error:
		return true
	default:
		return false
	}
}

func debugOutput() io.Writer {
	return os.Stdout
}
