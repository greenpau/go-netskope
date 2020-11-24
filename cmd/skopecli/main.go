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

package main

import (
	"flag"
	"fmt"
	"github.com/greenpau/go-netskope"
	"github.com/greenpau/versioned"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	app        *versioned.PackageManager
	appVersion string
	gitBranch  string
	gitCommit  string
	buildUser  string
	buildDate  string
)

func init() {
	app = versioned.NewPackageManager("skopecli")
	app.Description = "Netskope API Client"
	app.Documentation = "https://github.com/greenpau/go-netskope/"
	app.SetVersion(appVersion, "1.0.0")
	app.SetGitBranch(gitBranch, "main")
	app.SetGitCommit(gitCommit, "906996f")
	app.SetBuildUser(buildUser, "")
	app.SetBuildDate(buildDate, "")
}

func main() {
	var logLevel string
	var isShowVersion bool
	var configDir string
	var configFile string
	var token string
	var tenantName string
	var getClientData bool
	//var queryTimeRange uint64
	//var queryTimeOffset uint64

	flag.StringVar(&configFile, "config", "", "configuration file")
	flag.StringVar(&token, "token", "", "API Token")
	flag.StringVar(&tenantName, "tenant-name", "", "Tenant Name")

	flag.BoolVar(&getClientData, "get-client-data", false, "Get Client Data")

	//flag.Uint64Var(&queryTimeRange, "query-time-range", 0, "Time range in seconds, e.g. 3600, 86400, 604800, etc.")
	//flag.Uint64Var(&queryTimeOffset, "query-time-offset", 0, "Time offset in seconds")

	flag.StringVar(&logLevel, "log-level", "info", "logging severity level")
	flag.BoolVar(&isShowVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s - %s\n\n", app.Name, app.Description)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", app.Name)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation: %s\n\n", app.Documentation)
	}
	flag.Parse()

	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s\n", app.Banner())
		os.Exit(0)
	}

	// Determine configuration file name and extension
	if configFile == "" {
		configDir = "."
		configFile = app.Name + ".yaml"
	} else {
		configDir, configFile = filepath.Split(configFile)
	}
	configFileExt := filepath.Ext(configFile)
	if configFileExt == "" {
		fmt.Fprintf(os.Stderr, "--config specifies a file without an extension, e.g. .yaml or .json\n")
		os.Exit(1)
	}

	configName := strings.TrimSuffix(configFile, configFileExt)
	viper.SetConfigName(configName)
	viper.SetEnvPrefix("netskope")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_", " ", "_"))
	viper.AddConfigPath("$HOME/.config/" + app.Name)
	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	// Obtain settings via environment variable
	if token == "" {
		if v := viper.Get("token"); v != nil {
			token = viper.Get("token").(string)
		}
	}

	if tenantName == "" {
		if v := viper.Get("tenant-name"); v != nil {
			tenantName = viper.Get("tenant-name").(string)
		}
	}

	// Obtain settings via configuration file
	if err := viper.ReadInConfig(); err == nil {
		if token == "" {
			if v := viper.Get("token"); v != nil {
				token = viper.Get("token").(string)
			}
		}
		if tenantName == "" {
			if v := viper.Get("tenant_name"); v != nil {
				tenantName = viper.Get("tenant_name").(string)
			}
		}
	} else {
		if !strings.Contains(err.Error(), "Not Found in") {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	opts := make(map[string]interface{})
	opts["log_level"] = logLevel
	cli, err := netskope.NewClient(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer cli.Close()

	if err := cli.SetToken(token); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := cli.SetTenantName(tenantName); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	cli.Info()

	opts = make(map[string]interface{})
	if getClientData {
		items, err := cli.GetClientEndpoints(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for _, item := range items {
			s, err := item.ToJSONString()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "%s\n", s)
		}
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "actionable argument is missing\n")
	os.Exit(1)
}
