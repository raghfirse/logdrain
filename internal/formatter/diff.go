package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Diff compares consecutive JSON log lines and emits only changed fields.
type Diff struct {
	fields []string
	prev   map[string]json.RawMessage
}

// NewDiff creates a Diff processor. If fields is empty, all fields are compared.
func NewDiff(fields []string) *Diff {
	return &Diff{fields: fields}
}

// Apply compares the current line against the previous one and returns a line
// containing only the fields that changed (plus a special "_diff" marker).
// Non-JSON lines are passed through unchanged.
func (d *Diff) Apply(line string) string {
	var cur map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &cur); err != nil {
		d.prev = nil
		return line
	}

	if d.prev == nil {
		d.prev = cur
		return line
	}

	keys := d.fields
	if len(keys) == 0 {
		keys = allKeys(cur, d.prev)
	}

	changed := map[string]json.RawMessage{}
	for _, k := range keys {
		prevVal, hadPrev := d.prev[k]
		curVal, hasCur := cur[k]
		switch {
		case hasCur && hadPrev && string(curVal) != string(prevVal):
			changed[k] = curVal
		case hasCur && !hadPrev:
			changed[k] = curVal
		case !hasCur && hadPrev:
			changed[k] = json.RawMessage(`null`)
		}
	}

	d.prev = cur

	if len(changed) == 0 {
		changed["_diff"] = json.RawMessage(`"(no change)"`)
	} else {
		changed["_diff"] = json.RawMessage(fmt.Sprintf(`"%d field(s) changed"`, len(changed)))
	}

	out, err := json.Marshal(changed)
	if err != nil {
		return line
	}
	return string(out)
}

// Reset clears the previous line state.
func (d *Diff) Reset() {
	d.prev = nil
}

// ParseDiffFlag parses a comma-separated list of field names for diffing.
// An empty string returns an empty slice (diff all fields).
func ParseDiffFlag(s string) ([]string, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	var fields []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("diff: empty field name in %q", s)
		}
		fields = append(fields, p)
	}
	return fields, nil
}

func allKeys(a, b map[string]json.RawMessage) []string {
	seen := map[string]struct{}{}
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
