# go-netskope

<a href="https://github.com/greenpau/go-netskope/actions/" target="_blank"><img src="https://github.com/greenpau/go-netskope/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/go-netskope" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>

Netskope API Client Library

<!-- begin-markdown-toc -->
## Table of Contents

* [Getting Started](#getting-started)
* [References](#references)

<!-- end-markdown-toc -->

## Getting Started

First, install `skopecli`:

```bash
go get -u github.com/greenpau/go-netskope/cmd/skopecli
```

Next, set environment variables for Netskope API Token:

```bash
export NETSKOPE_TOKEN=8c5e8866-0062-4059-b2be-92707e4374da
export NETSKOPE_TENANT_NAME=acmeprod
```

Alternatively, the settings could be passed in a configuration file. There are
two options:

1. The `skopecli.yaml` should be located in `$HOME/.config/skopecli` or current directory
2. Pass the location via `-config` flag

```yaml
---
token: 8c5e8866-0062-4059-b2be-92707e4374da
tenant_name: acmeprod
```

The following command fetches client data from Netskope API:

```bash
skopecli -get-client-data
```
