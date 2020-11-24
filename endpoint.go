// Copyright 2020 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netskope

import (
	"encoding/json"
	"fmt"
	//"go.uber.org/zap"
	"strconv"
)

// ClientEndpointResponse is response from clients API endpoint.
type ClientEndpointResponse struct {
	Status          string            `json:"success,omitempty"`
	Message         string            `json:"msg,omitempty"`
	ClientEndpoints []*ClientEndpoint `json:"data,omitempty"`
}

// ClientEndpoint is a compute endpoint.
type ClientEndpoint struct {
	ID               string  `json:"_id,omitempty"`
	DeviceID         string  `json:"device_id,omitempty"`
	InstallTimestamp float64 `json:"client_install_time,omitempty"`
	Version          string  `json:"client_version,omitempty"`
	Users            []*User `json:"users,omitempty"`
	HostInfo         *Host   `json:"host_info,omitempty"`
	LastEvent        *Event  `json:"last_event,omitempty"`
}

// GetClientEndpoints returns a list of ClientEndpoint instances.
func (c *Client) GetClientEndpoints(opts map[string]interface{}) ([]*ClientEndpoint, error) {
	endpoints := []*ClientEndpoint{}
	pageSize := 1000
	offset := 0
	if v, exists := opts["page_size"]; exists {
		pageSize = v.(int)
	}
	for {
		params := make(map[string]string)
		params["limit"] = strconv.Itoa(pageSize)
		params["skip"] = strconv.Itoa(offset)
		b, err := c.callAPI("GET", "clients", params)
		if err != nil {
			return endpoints, err
		}
		resp := &ClientEndpointResponse{}
		if err := json.Unmarshal(b, &resp); err != nil {
			return endpoints, fmt.Errorf("failed unmarshalling response: %s", err)
		}
		//c.log.Warn("response", zap.Any("response", string(b)))
		//if resp.Status != "success" {
		//	return endpoints, fmt.Errorf("failed request: %s", resp.Message)
		//}
		for _, endpoint := range resp.ClientEndpoints {
			//c.log.Warn("endpoint", zap.Any("endpoint", endpoint))
			endpoints = append(endpoints, endpoint)
		}

		offset += len(resp.ClientEndpoints)
		if len(resp.ClientEndpoints) < pageSize {
			break
		}
	}
	return endpoints, nil
}

// ToJSONString serializes ClientEndpoint to a string.
func (c *ClientEndpoint) ToJSONString() (string, error) {
	itemJSON, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed converting to json: %s", err)
	}
	return string(itemJSON), nil
}

// UnmarshalJSON unpacks byte array into ClientEndpoint.
func (c *ClientEndpoint) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if len(b) < 10 {
		return fmt.Errorf("invalid ClientEndpoint data: %s", b)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unpack ClientEndpoint")
	}

	if _, exists := m["attributes"]; !exists {
		return fmt.Errorf("failed to unpack ClientEndpoint, attributes not found")
	}

	for k, v := range m["attributes"].(map[string]interface{}) {
		switch k {
		case "_id":
			c.ID = v.(string)
		case "device_id":
			c.DeviceID = v.(string)
		case "client_install_time":
			c.InstallTimestamp = v.(float64)
		case "client_version":
			c.Version = v.(string)
		case "users":
			for _, u := range v.([]interface{}) {
				usr := &User{}
				if err := usr.load(u.(map[string]interface{})); err != nil {
					return fmt.Errorf("failed to unpack ClientEndpoint, %s attribute error: %s", k, err)
				}
				c.Users = append(c.Users, usr)
			}
		case "host_info":
			hostInfo := &Host{}
			if err := hostInfo.load(v.(map[string]interface{})); err != nil {
				return fmt.Errorf("failed to unpack ClientEndpoint, %s attribute error: %s", k, err)
			}
			c.HostInfo = hostInfo
		case "last_event":
			lastEvent := &Event{}
			if err := lastEvent.load(v.(map[string]interface{})); err != nil {
				return fmt.Errorf("failed to unpack ClientEndpoint, %s attribute error: %s", k, err)
			}
			c.LastEvent = lastEvent
		default:
			return fmt.Errorf("failed to unpack ClientEndpoint, unsupported attribute: %s, %v", k, v)
		}
	}

	return nil
}
