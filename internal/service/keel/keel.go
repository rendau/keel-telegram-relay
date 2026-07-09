package keel

import (
	"encoding/json"
	"fmt"
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
//
// The first line leads with the affected resource name so it is visible in the
// notification preview (a locked screen shows only the first line). Everything
// that is not essential (redundant name, verbose metadata, nanosecond
// timestamp) is dropped — telegram already shows the message time.
func (e *Event) Text() string {
	var sb strings.Builder

	sb.WriteString(levelIcon(e.Level))
	if name := e.resourceName(); name != "" {
		sb.WriteString(" " + name)
	} else {
		sb.WriteString(" Keel")
	}
	if e.Type != "" {
		sb.WriteString(" · " + e.Type)
	}

	if e.Message != "" {
		sb.WriteString("\n" + e.Message)
	}

	return strings.TrimRight(sb.String(), "\n")
}

// resourceName returns the name of the affected resource (e.g. the deployment
// name) taken from metadata, or "" if it is not present.
func (e *Event) resourceName() string {
	v, ok := e.Metadata["name"]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
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
