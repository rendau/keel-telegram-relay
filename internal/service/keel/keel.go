package keel

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Event is a keel notification payload sent to a generic webhook endpoint.
//
// Example (deployment update):
//
//	{
//	  "name": "update deployment",
//	  "message": "Successfully updated deployment default/oms (ghcr.io/mechta-market/oms:1.2.3)",
//	  "type": "deployment update",
//	  "level": "success",
//	  "createdAt": "2026-07-08T10:00:00Z",
//	  "metadata": {"provider": "kubernetes", "namespace": "default"}
//	}
type Event struct {
	Name      string         `json:"name"`
	Message   string         `json:"message"`
	Type      string         `json:"type"`
	Level     string         `json:"level"`
	CreatedAt string         `json:"createdAt"`
	Metadata  map[string]any `json:"metadata"`
}

func Parse(data []byte) (*Event, error) {
	event := &Event{}

	if err := json.Unmarshal(data, event); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return event, nil
}

// Text builds a plain-text telegram message (no parse_mode, so nothing to escape).
func (e *Event) Text() string {
	var sb strings.Builder

	sb.WriteString(levelIcon(e.Level) + " Keel")
	if e.Type != "" {
		sb.WriteString(" · " + e.Type)
	}
	sb.WriteString("\n")

	if e.Name != "" {
		sb.WriteString(e.Name + "\n")
	}
	if e.Message != "" {
		sb.WriteString("\n" + e.Message + "\n")
	}

	if len(e.Metadata) > 0 {
		sb.WriteString("\n")
		for _, k := range sortedKeys(e.Metadata) {
			sb.WriteString(fmt.Sprintf("%s: %v\n", k, e.Metadata[k]))
		}
	}

	if e.CreatedAt != "" {
		sb.WriteString("\n" + e.CreatedAt)
	}

	return strings.TrimRight(sb.String(), "\n")
}

func levelIcon(level string) string {
	switch strings.ToLower(level) {
	case "success":
		return "✅"
	case "error", "failure", "fatal":
		return "🔴"
	case "warn", "warning":
		return "⚠️"
	default:
		return "🚀"
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
