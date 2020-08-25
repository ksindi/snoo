# SNOO

![GitHub last commit](https://img.shields.io/github/last-commit/ksindi/snoo?style=for-the-badge)

An API client to the [SNOO Smart Sleeper Bassinet](https://www.happiestbaby.com/products/snoo-smart-bassinet).
The SNOO is a bassinet that will rock your baby to sleep, and responds to the
baby by trying to sooth it with different rocking motions and sounds when it
detects crying.

Currently, it supports getting the current session data from SNOO, and historic
data. It does not allow you to control the SNOO (the control API is provided by
[PubNub](https://www.pubnub.com) and is different from the read-only data API
hosted by happiestbaby.com)

# Disclaimer

The SNOO API is undocumented. Using it might or might not violate Happiest Baby, Inc
[Terms of Service](https://www.happiestbaby.com/pages/terms-of-service).
Use at your own risk.

## Install

You can also just grab the latest version with `curl`. For Linux:

```sh
sudo curl -o /usr/local/bin/snoo https://github.com/ksindi/snoo/releases/download/v0.1.0/snoo-linux
sudo chmod +x /usr/local/bin/snoo
```

Or on MacOS:

```sh
sudo curl -o /usr/local/bin/snoo https://github.com/ksindi/snoo/releases/download/v0.1.0/snoo-darwin
sudo chmod +x /usr/local/bin/snoo
```

## Usage

```sh
$ snoo -h

NAME:
   snoo - An API client to the SNOO Smart Sleeper Bassinet

USAGE:
   snoo [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   status, s  get current status of SNOO
   sessions   print list of sessions for a date range
   days       print list of daily aggregated sessions for a date range
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value  SNOO login username [$SNOO_USERNAME]
   --password value  SNOO login password [$SNOO_PASSWORD]
   --help, -h        show help (default: false)
   --version, -v     print the version (default: false)
```

### Get status

To get the status of your snoo, simply run

```sh
$ snoo status
Soothing 26m
```

The output of the `snoo status` command is the status
(`Awake`, `Asleep`, or `Soothing`), and the duration of the current session.

### Export daily aggregated data

To export daily aggregated data, use

```sh
$ snoo days --start DATE --end DATE
```

```csv
date,naps,longest_sleep,total_sleep,day_sleep,night_sleep,night_wakings,timezone
2019-12-03,6,12933,58035,25038,32997,4,None
```

Again, all durations are given in seconds. How `day_sleep` and `night_sleep` are defined is set in your Snoo app.

## Programmatic usage

```go
package main

import (
  "context"
  "log"
  "net/http"

  "github.com/ksindi/snoo"
)

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  client := snoo.NewClient("my_username", "my_password")

  // declare an empty interface
  var result map[string]interface{}
  err := client.MakeRequest(ctx, http.MethodGet, "/ss/v2/sessions/last", nil, nil, &result)

  if err != nil {
  	log.Fatal(err)
  }
}
```

## Credit

This repo was inspired by https://github.com/maebert/snoo.
