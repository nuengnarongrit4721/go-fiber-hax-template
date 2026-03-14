package logs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type PrettyHandler struct {
	w      io.Writer
	opts   slog.HandlerOptions
	attrs  []slog.Attr
	groups []string
	mu     sync.Mutex
}

const timeFormat = "2006-01-02 15:04:05.000 -07:00"

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	return &PrettyHandler{w: w, opts: *opts}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}

	header := fmt.Sprintf("[%s] [%s] %s\n", ts.Format(timeFormat), strings.ToUpper(r.Level.String()), r.Message)
	if _, err := io.WriteString(h.w, header); err != nil {
		return err
	}

	attrs := make([]slog.Attr, 0, len(h.attrs)+r.NumAttrs())
	attrs = append(attrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	prefix := strings.Join(h.groups, ".")
	for _, a := range attrs {
		h.writeAttr(prefix, a)
	}

	_, err := io.WriteString(h.w, "\n")
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cloned := h.clone()
	cloned.attrs = append(cloned.attrs, attrs...)
	return cloned
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	cloned := h.clone()
	cloned.groups = append(cloned.groups, name)
	return cloned
}

func (h *PrettyHandler) clone() *PrettyHandler {
	return &PrettyHandler{
		w:      h.w,
		opts:   h.opts,
		attrs:  append([]slog.Attr(nil), h.attrs...),
		groups: append([]string(nil), h.groups...),
	}
}

func (h *PrettyHandler) writeAttr(prefix string, a slog.Attr) {
	a = slog.Attr{Key: a.Key, Value: a.Value.Resolve()}
	key := a.Key
	if prefix != "" {
		key = prefix + "." + key
	}

	v := a.Value
	if v.Kind() == slog.KindGroup {
		for _, ga := range v.Group() {
			h.writeAttr(key, ga)
		}
		return
	}

	if pv, ok := v.Any().(PrettyValue); ok {
		pretty := pv.String()
		io.WriteString(h.w, fmt.Sprintf("  %s:\n", key))
		io.WriteString(h.w, indentLines(pretty, "    "))
		io.WriteString(h.w, "\n")
		return
	}

	io.WriteString(h.w, fmt.Sprintf("  %s: %s\n", key, valueString(v)))
}

func valueString(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%f", v.Float64())
	case slog.KindBool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339Nano)
	case slog.KindAny:
		return fmt.Sprint(v.Any())
	default:
		return v.String()
	}
}

func indentLines(s, prefix string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}
