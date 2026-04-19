package runtimelog

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// Handler is a slog handler that fans out records into:
//   - the wrapped next handler
//   - the in-memory Store for TUI consumption
//
// The handler expects runtime loggers to include bot_id via logger.With(...).
type Handler struct {
	next   slog.Handler
	store  *Store
	attrs  []slog.Attr
	groups []string
}

// NewHandler constructs a runtime log fan-out handler.
func NewHandler(next slog.Handler, store *Store) *Handler {
	return &Handler{
		next:  next,
		store: store,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	if h == nil {
		return false
	}
	if h.next != nil {
		return h.next.Enabled(ctx, level)
	}
	return true
}

// Handle writes one record to the in-memory store and then to the wrapped handler.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	if h == nil {
		return nil
	}

	entry := Entry{
		Time:    record.Time,
		Level:   record.Level,
		Message: record.Message,
		Fields:  make([]Field, 0, len(h.attrs)+record.NumAttrs()),
	}

	for _, attr := range h.attrs {
		appendAttrToEntry(h.groups, attr, &entry)
	}

	record.Attrs(func(attr slog.Attr) bool {
		appendAttrToEntry(h.groups, attr, &entry)
		return true
	})

	if h.store != nil {
		h.store.Append(entry)
	}

	if h.next != nil {
		return h.next.Handle(ctx, record)
	}

	return nil
}

// WithAttrs returns a new handler with appended attributes.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h == nil {
		return nil
	}

	var next slog.Handler
	if h.next != nil {
		next = h.next.WithAttrs(attrs)
	}

	combined := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	combined = append(combined, h.attrs...)
	combined = append(combined, attrs...)

	groups := make([]string, len(h.groups))
	copy(groups, h.groups)

	return &Handler{
		next:   next,
		store:  h.store,
		attrs:  combined,
		groups: groups,
	}
}

// WithGroup returns a new handler under one more group.
func (h *Handler) WithGroup(name string) slog.Handler {
	if h == nil {
		return nil
	}

	var next slog.Handler
	if h.next != nil {
		next = h.next.WithGroup(name)
	}

	attrs := make([]slog.Attr, len(h.attrs))
	copy(attrs, h.attrs)

	groups := make([]string, 0, len(h.groups)+1)
	groups = append(groups, h.groups...)
	if strings.TrimSpace(name) != "" {
		groups = append(groups, name)
	}

	return &Handler{
		next:   next,
		store:  h.store,
		attrs:  attrs,
		groups: groups,
	}
}

func appendAttrToEntry(groups []string, attr slog.Attr, entry *Entry) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return
	}

	key := joinKey(groups, attr.Key)

	if attr.Value.Kind() == slog.KindGroup {
		childGroups := groups
		if strings.TrimSpace(attr.Key) != "" {
			childGroups = appendGroup(groups, attr.Key)
		}

		for _, child := range attr.Value.Group() {
			appendAttrToEntry(childGroups, child, entry)
		}
		return
	}

	value := attrValueString(attr.Value)

	entry.Fields = append(entry.Fields, Field{
		Key:   key,
		Value: value,
	})

	switch key {
	case "bot_id":
		entry.BotID = value
	case "bot_name":
		entry.BotName = value
	}
}

func appendGroup(groups []string, name string) []string {
	out := make([]string, 0, len(groups)+1)
	out = append(out, groups...)
	out = append(out, name)
	return out
}

func joinKey(groups []string, key string) string {
	switch {
	case len(groups) == 0:
		return key
	case key == "":
		return strings.Join(groups, ".")
	default:
		return strings.Join(appendGroup(groups, key), ".")
	}
}

func attrValueString(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindBool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%v", v.Float64())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format("2006-01-02 15:04:05")
	case slog.KindAny:
		return fmt.Sprint(v.Any())
	default:
		return fmt.Sprint(v)
	}
}
