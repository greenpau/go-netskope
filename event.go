package netskope

import (
	"fmt"
)

// Event is an event and its attributes.
type Event struct {
	Actor     string  `json:"actor,omitempty"`
	Name      string  `json:"event,omitempty"`
	NpaStatus string  `json:"npa_status,omitempty"`
	Status    string  `json:"status,omitempty"`
	Timestamp float64 `json:"timestamp,omitempty"`
}

func (e *Event) load(m map[string]interface{}) error {
	if m == nil {
		return nil
	}
	for k, v := range m {
		switch k {
		case "actor":
			e.Actor = v.(string)
		case "event":
			e.Name = v.(string)
		case "npa_status":
			e.NpaStatus = v.(string)
		case "status":
			e.Status = v.(string)
		case "timestamp":
			e.Timestamp = v.(float64)
		default:
			return fmt.Errorf("unsupported attribute: %s, %v", k, v)
		}
	}
	return nil
}
