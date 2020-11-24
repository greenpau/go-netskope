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
	"crypto/tls"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ReceiverDataLimit is the limit of data in bytes the client will read
// from a server.
const ReceiverDataLimit int64 = 1e6

// Client is an instance of Proofpoint API client.
type Client struct {
	url                string
	host               string
	port               int
	protocol           string
	token              string
	tenantName         string
	validateServerCert bool
	dataLimit          int64
	pathPrefix         string
	log                *zap.Logger
}

// NewClient returns an instance of Client.
func NewClient(opts map[string]interface{}) (*Client, error) {
	c := &Client{
		host:       "tenant.goskope.com",
		port:       443,
		protocol:   "https",
		pathPrefix: "/api/v1/",
		dataLimit:  ReceiverDataLimit,
	}
	log, err := newLogger(opts)
	if err != nil {
		return nil, fmt.Errorf("failed initializing log: %s", err)
	}
	c.log = log
	return c, nil
}

// Close performs a cleanup associated with Client..
func (c *Client) Close() {
	if c.log != nil {
		c.log.Sync()
	}
}

// Info sends information about Client to the configured logger.
func (c *Client) Info() {
	c.rebaseURL()
	c.log.Debug(
		"client configuration",
		zap.String("url", c.url),
		zap.String("path_prefix", c.pathPrefix),
	)
}

func (c *Client) rebaseURL() {
	if (c.protocol == "https" && c.port == 443) ||
		(c.protocol == "http" && c.port == 80) {
		c.url = fmt.Sprintf("%s://%s", c.protocol, c.host)
		return
	}
	c.url = fmt.Sprintf("%s://%s:%d", c.protocol, c.host, c.port)
	return
}

// SetHost sets the target host for the API calls.
func (c *Client) SetHost(s string) error {
	if s == "" {
		return fmt.Errorf("empty hostname or ip address")
	}
	c.host = s
	c.rebaseURL()
	return nil
}

// SetPort sets the port number for the API calls.
func (c *Client) SetPort(p int) error {
	if p == 0 {
		return fmt.Errorf("invalid port: %d", p)
	}
	c.port = p
	c.rebaseURL()
	return nil
}

// SetToken sets API token.
func (c *Client) SetToken(s string) error {
	if s == "" {
		return fmt.Errorf("empty token")
	}
	c.token = s
	return nil
}

// SetTenantName sets API tenant name.
func (c *Client) SetTenantName(s string) error {
	if s == "" {
		return fmt.Errorf("empty tenant name")
	}
	c.tenantName = s
	c.host = strings.ReplaceAll(c.host, "tenant", c.tenantName)
	return c.SetHost(c.host)
}

// SetProtocol sets the protocol for the API calls.
func (c *Client) SetProtocol(s string) error {
	switch s {
	case "http":
		c.protocol = s
	case "https":
		c.protocol = s
	default:
		return fmt.Errorf("supported protocols: http, https; unsupported protocol: %s", s)
	}
	c.rebaseURL()
	return nil
}

// SetValidateServerCertificate instructs the client to enforce the validation of certificates
// and check certificate errors.
func (c *Client) SetValidateServerCertificate() error {
	c.validateServerCert = true
	return nil
}

func (c *Client) callAPI(method, svc string, params map[string]string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s%s%s", c.url, c.pathPrefix, svc)
	c.log.Debug(
		"making http request",
		zap.String("method", method),
		zap.String("url", reqURL),
		zap.Any("params", params),
	)

	q := url.Values{}
	q.Set("token", c.token)
	for k, v := range params {
		q.Set(k, v)
	}

	reqURL = fmt.Sprintf("%s?%s", reqURL, q.Encode())

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	if !c.validateServerCert {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}

	var req *http.Request
	var err error
	req, err = http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, err
	}

	//if contentType != "" {
	//	req.Header.Add("Content-Type", contentType)
	//}
	req.Header.Add("Accept", "application/json;charset=utf-8")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := httpClient.Do(req)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "EOF") {
			return nil, err
		}
	}
	if res == nil {
		return nil, fmt.Errorf("response: <nil>, verify url: %s", reqURL)
	}
	defer res.Body.Close()

	c.log.Debug("http response", zap.String("status", res.Status))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("non-EOF error at url %s: %s", reqURL, err)
	}

	// c.log.Debug("http response body", zap.String("body", string(body)))

	switch res.StatusCode {
	case 200:
		return body, nil
	case 204:
		return body, nil
	default:
		return nil, fmt.Errorf("error: status code %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}
