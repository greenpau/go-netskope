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
	"fmt"
	. "github.com/greenpau/go-netskope/internal/server"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	timerStartTime := time.Now()
	token := "8c5e8866-0062-4059-b2be-92707e4374da"

	// Create a client instance
	opts := make(map[string]interface{})
	opts["log_level"] = "debug"
	cli, err := NewClient(opts)
	if err != nil {
		t.Fatalf("failed initializing client: %s", err)
	}
	defer cli.Close()

	// Create web server instance
	endpoints := map[string][]*MockTestEndpoint{
		"/api/v1/clients": []*MockTestEndpoint{
			&MockTestEndpoint{
				RequestURI: fmt.Sprintf("/api/v1/clients?token=%s", token),
				FileName:   "clients.json",
			},
		},
	}

	server, err := NewMockTestServer(cli.log, endpoints, token, true)
	if err != nil {
		t.Fatalf("Failed to initialize mock test server: %s", err)
	}
	defer server.Close()

	// Configure client
	cli.SetHost(server.NonTLS.Hostname)
	cli.SetPort(server.NonTLS.Port)
	cli.SetProtocol(server.NonTLS.Protocol)
	cli.SetToken(token)
	if err := cli.SetValidateServerCertificate(); err != nil {
		t.Fatalf("expected success, but failed")
	}
	cli.Info()

	t.Logf("client: took %s", time.Since(timerStartTime))
}
