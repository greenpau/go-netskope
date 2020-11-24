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

package server

import (
	"fmt"

	"go.uber.org/zap"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

// MockTestServerInstance is an instance of a mock web server.
type MockTestServerInstance struct {
	Instance *httptest.Server
	URL      *url.URL
	Hostname string
	Protocol string
	Port     int
}

// MockTestServer is a mock web server. The server supports both HTTPS and HTTP.
type MockTestServer struct {
	NonTLS *MockTestServerInstance
	TLS    *MockTestServerInstance
}

// MockTestEndpoint is a mock API endpoint.
type MockTestEndpoint struct {
	RequestURI string
	FileName   string
}

// Close closes running instances of MockTestServerInstance, if any.
func (srv *MockTestServer) Close() {
	if srv.NonTLS.Instance != nil {
		srv.NonTLS.Instance.Close()
	}
	if srv.TLS.Instance != nil {
		srv.TLS.Instance.Close()
	}
}

// NewMockTestServer return an instance of MockTestServer running
// with and without TLS.
func NewMockTestServer(log *zap.Logger, pathMap map[string][]*MockTestEndpoint, token string, tlsEnabled bool) (*MockTestServer, error) {
	// Create web server instance
	mts := &MockTestServer{
		NonTLS: &MockTestServerInstance{},
		TLS:    &MockTestServerInstance{},
	}
	serverEndpoints := map[string][]*MockTestEndpoint{
		"/v2/siem/all": []*MockTestEndpoint{
			&MockTestEndpoint{
				RequestURI: "/v2/siem/all?format=json&sinceSeconds=3600",
				FileName:   "siem_all_format_json_since_seconds_3600.json",
			},
		},
		"/v2/siem/issues": []*MockTestEndpoint{
			&MockTestEndpoint{
				RequestURI: "/v2/siem/all?format=json&sinceSeconds=3600",
				FileName:   "siem_issues_format_json_since_seconds_3600.json",
			},
		},
	}

	if pathMap != nil {
		for k, v := range pathMap {
			if _, exists := serverEndpoints[k]; !exists {
				serverEndpoints[k] = v
			} else {
				for _, endpoint := range v {
					serverEndpoints[k] = append(serverEndpoints[k], endpoint)
				}
			}
		}
	}

	log.Debug("api endpoints", zap.Any("endpoints", serverEndpoints))

	dataDir := "testdata/responses"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var err error
		var fp string
		var fc []byte
		isAuthError := true

		// TODO: add token validation
		isAuthError = false

		if isAuthError {
			fp = fmt.Sprintf("%s/access_denied_error_1.json", dataDir)
			fc, err = ioutil.ReadFile(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(fc), http.StatusNotFound)
			return
		}

		if req.Method != "GET" {
			http.Error(w, "Bad Request, expecting GET", http.StatusBadRequest)
			return
		}

		if strings.HasSuffix(req.URL.Path, "/empty_response") {
			panic("")
		}

		if strings.HasSuffix(req.URL.Path, "/replay_request") {
			reqBlob, _ := httputil.DumpRequest(req, true)
			w.Write(reqBlob)
			return
		}

		// TODO: endpoints with various options
		endpoints, respFileExists := serverEndpoints[req.URL.Path]
		if !respFileExists {
			fp = fmt.Sprintf("%s/not_found_error_1.json", dataDir)
			fc, err = ioutil.ReadFile(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(fc), http.StatusNotFound)
			return
		}
		fp = fmt.Sprintf("%s/%s", dataDir, endpoints[0].FileName)
		fc, err = ioutil.ReadFile(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(fc)
	})

	// Initialize HTTP server
	mts.NonTLS.Instance = httptest.NewServer(mux)

	httpServerURL, err := url.Parse(mts.NonTLS.Instance.URL)
	if err != nil {
		return nil, err
	}
	mts.NonTLS.URL = httpServerURL
	log.Debug("configured url", zap.String("url", mts.NonTLS.Instance.URL))

	mts.NonTLS.Hostname = httpServerURL.Hostname()
	log.Debug("configured hostname", zap.String("hostname", mts.NonTLS.Hostname))

	mts.NonTLS.Protocol = strings.Split(mts.NonTLS.Instance.URL, ":")[0]
	log.Debug("configured protocol", zap.String("protocol", mts.NonTLS.Protocol))

	httpServerPort, _ := strconv.Atoi(httpServerURL.Port())
	mts.NonTLS.Port = httpServerPort
	log.Debug("configured port", zap.Int("port", mts.NonTLS.Port))

	if tlsEnabled {
		// Initialize HTTPS server
		mts.TLS.Instance = httptest.NewTLSServer(mux)
		httpsServerURL, err := url.Parse(mts.TLS.Instance.URL)
		if err != nil {
			return nil, err
		}
		mts.TLS.URL = httpsServerURL
		mts.TLS.Hostname = httpsServerURL.Hostname()
		log.Debug("configured hostname", zap.String("hostname", mts.TLS.Hostname))

		mts.TLS.Protocol = strings.Split(mts.TLS.Instance.URL, ":")[0]
		log.Debug("configured protocol", zap.String("protocol", mts.TLS.Protocol))

		httpsServerPort, _ := strconv.Atoi(httpsServerURL.Port())
		mts.TLS.Port = httpsServerPort
		log.Debug("configured port", zap.Int("port", mts.TLS.Port))
	}

	return mts, nil
}
