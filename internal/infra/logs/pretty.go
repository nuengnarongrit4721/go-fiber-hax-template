package logs

import "encoding/json"

type PrettyValue struct {
	V any
}

func Pretty(v any) PrettyValue {
	return PrettyValue{V: v}
}

func (p PrettyValue) String() string {
	b, err := json.MarshalIndent(p.V, "", "  ")
	if err != nil {
		return "pretty_error"
	}
	return string(b)
}
