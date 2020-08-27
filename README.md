# SNOO

![GitHub last commit](https://img.shields.io/github/last-commit/ksindi/snoo?style=for-the-badge)

An API client to the [SNOO Smart Sleeper Bassinet](https://www.happiestbaby.com/products/snoo-smart-bassinet).
The SNOO is a bassinet that will rock your baby to sleep. It responds to
baby by trying to sooth it with different rocking motions and sounds when it
detects crying.

The client only supports getting session and daily aggregated data
from the SNOO. It does not allow you to control the SNOO. That API is provided
by [PubNub](https://www.pubnub.com) and is different from the read-only data
API hosted by Happiest Baby.

## Disclaimer

The SNOO API is undocumented. Using it might or might not violate Happiest Baby, Inc's
[Terms of Service](https://www.happiestbaby.com/pages/terms-of-service).
Use at your own risk.

## Install

You can grab the latest version with `curl`. For Linux:

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
   sessions   print list of sessions for a date range
   days       print list of daily aggregated sessions for a date range
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value  SNOO login username [$SNOO_USERNAME]
   --password value  SNOO login password [$SNOO_PASSWORD]
   --debug           Enable verbose debugging (default: false)
   --help, -h        show help (default: false)
   --version, -v     print the version (default: false)
```

### Export daily aggregated data

To export daily aggregated data, use

```sh
$ snoo days --start DATE --end DATE
```

```csv
date,naps,longest_sleep,total_sleep,day_sleep,night_sleep,night_wakings,timezone
2020-08-01 00:00:00 +0000 UTC,4,12997,55829,33950,21879,3,
```

All durations are given in seconds. How `day_sleep` and `night_sleep`
are defined is set in your Snoo app.

### Export sessions

To export sessions, use

```sh
$ snoo sessions --start DATE --end DATE
```

```csv
session_id,start_time,end_time,duration,asleep_duration,soothing_duration
1572833733,2020-08-01 00:00:00,2020-08-01 02:55:40,10541,9682,85
```

Again, all durations are given in seconds. How `asleep_duration` and `soothing_duration`
are defined is set in your Snoo app.

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
  	log.Fatalln(err)
  }
}
```

## Credit

This repo is inspired by https://github.com/maebert/snoo.
