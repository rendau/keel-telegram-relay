package keel

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
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
// Three lines:
//
//	✅ <name>
//	repository: <repo>
//	<time>
//
// The first line leads with a level icon and the affected resource name so it
// is visible in the notification preview (a locked screen shows only the first
// line). The second line carries the image repository, the third the event time.
func (e *Event) Text() string {
	var sb strings.Builder

	sb.WriteString(levelIcon(e.Level))
	if name := e.resourceName(); name != "" {
		sb.WriteString(" " + name)
	} else {
		sb.WriteString(" Keel")
	}

	if repo := e.repository(); repo != "" {
		sb.WriteString("\nrepository: " + repo)
	}

	if t := e.eventTime(); t != "" {
		sb.WriteString("\n" + t)
	}

	return strings.TrimRight(sb.String(), "\n")
}

// resourceName returns the name of the affected resource, preferring the
// namespace/name from metadata (e.g. "default/oms") and falling back to the
// bare name or the event name.
func (e *Event) resourceName() string {
	name := e.metaString("name")
	if name == "" {
		return strings.TrimSpace(e.Name)
	}
	if ns := e.metaString("namespace"); ns != "" {
		return ns + "/" + name
	}
	return name
}

// parenRe captures the last parenthesised group of the keel message, which is
// where the image reference lives, e.g. "... (ghcr.io/mechta-market/oms:1.2.3)".
var parenRe = regexp.MustCompile(`\(([^)]+)\)`)

// imageRe matches a docker image reference (a slashed path optionally followed
// by a ":tag") anywhere in a string.
var imageRe = regexp.MustCompile(`[\w.\-]+(?:/[\w.\-]+)+(?::[\w.\-]+)?`)

// repository returns the image repository, taken from metadata when present or
// parsed out of the message text (preferring the parenthesised image ref).
func (e *Event) repository() string {
	if repo := e.metaString("repository"); repo != "" {
		return repo
	}
	if m := parenRe.FindStringSubmatch(e.Message); m != nil {
		if img := imageRe.FindString(m[1]); img != "" {
			return img
		}
	}
	return imageRe.FindString(e.Message)
}

// eventTime returns the event time, normalised to "2006-01-02 15:04:05" when it
// parses as an RFC3339 timestamp, or the raw value otherwise.
func (e *Event) eventTime() string {
	raw := strings.TrimSpace(e.CreatedAt)
	if raw == "" {
		return ""
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return raw
}

// metaString returns the trimmed string value for key in metadata, or "".
func (e *Event) metaString(key string) string {
	v, ok := e.Metadata[key]
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
