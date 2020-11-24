package netskope

import (
	"fmt"
)

// Host is a host with its attributes.
type Host struct {
	DeviceMake             string `json:"device_make,omitempty"`
	DeviceModel            string `json:"device_model,omitempty"`
	Hostname               string `json:"hostname,omitempty"`
	ManagementID           string `json:"managementID,omitempty"`
	OperatingSystemName    string `json:"os,omitempty"`
	OperatingSystemVersion string `json:"os_version,omitempty"`
	NsDeviceUID            string `json:"nsdeviceuid,omitempty"`
}

func (h *Host) load(m map[string]interface{}) error {
	if m == nil {
		return nil
	}

	for k, v := range m {
		switch k {
		case "device_make":
			h.DeviceMake = v.(string)
		case "device_model":
			h.DeviceModel = v.(string)
		case "hostname":
			h.Hostname = v.(string)
		case "managementID":
			h.ManagementID = v.(string)
		case "os":
			h.OperatingSystemName = v.(string)
		case "os_version":
			h.OperatingSystemVersion = v.(string)
		case "nsdeviceuid":
			h.NsDeviceUID = v.(string)
		default:
			return fmt.Errorf("unsupported attribute: %s, %v", k, v)
		}
	}
	return nil
}
