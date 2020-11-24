package netskope

import (
	"fmt"
)

// User is a user with its attributes.
type User struct {
	ID                         string `json:"_id,omitempty"`
	DeviceClassificationStatus string `json:"device_classification_status,omitempty"`
	LastEvent                  *Event
	OrganizationUnit           string   `json:"organization_unit,omitempty"`
	AddedTimestamp             float64  `json:"user_added_time,omitempty"`
	Groups                     []string `json:"user_groups,omitempty"`
	Source                     string   `json:"user_source,omitempty"`
	Key                        string   `json:"userkey,omitempty"`
	Username                   string   `json:"username,omitempty"`
}

func (u *User) load(m map[string]interface{}) error {
	if m == nil {
		return nil
	}
	for k, v := range m {
		switch k {
		case "_id":
			u.ID = v.(string)
		case "device_classification_status":
			u.DeviceClassificationStatus = v.(string)
		case "last_event":
			e := &Event{}
			if err := e.load(v.(map[string]interface{})); err != nil {
				return fmt.Errorf("failed to unpack ClientEndpoint, %s attribute error: %s", k, err)
			}
			u.LastEvent = e
		case "organization_unit":
			u.OrganizationUnit = v.(string)
		case "user_added_time":
			u.AddedTimestamp = v.(float64)
		case "user_groups":
			for _, g := range v.([]interface{}) {
				u.Groups = append(u.Groups, g.(string))
			}
		case "user_source":
			u.Source = v.(string)
		case "userkey":
			u.Key = v.(string)
		case "username":
			u.Username = v.(string)
		default:
			return fmt.Errorf("unsupported attribute: %s, %v", k, v)
		}
	}
	return nil
}
