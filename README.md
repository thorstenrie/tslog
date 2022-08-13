# tslog

[![Go Report Card](https://goreportcard.com/badge/github.com/thorstenrie/tslog)](https://goreportcard.com/report/github.com/thorstenrie/tslog)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/thorstenrie/tslog)](https://pkg.go.dev/mod/github.com/thorstenrie/tslog)

[Go](https://go.dev/) package for logging that tries to keep it simple ([KISS principle](https://en.wikipedia.org/wiki/KISS_principle)).

- **Simple**: Configured with one environment variable, initialized at startup of your app
- **Flexible**: Logging can be configured to stdout (default), to a temp file, a specifically defined file or even discarded.
- **Tested**: Unit tests including fuzzing with high [code coverage](https://gocover.io/github.com/thorstenrie/tslog)
- **Dependencies**: Only depends on [Go Standard Library](https://pkg.go.dev/std)

## Usage

Before app execution, set the environment variable `TS_LOGFILE`to your logging target (see [Configuration](#Configuration)).

E.g., in a linux terminal run

```
export TS_LOGFILE=log.txt
```

E.g., in VS Code add to the `configuration` block:
```
"env": {
    "TS_LOGFILE": "log.txt"
}
```

In the Go app, the package is imported with

```
import "github.com/thorstenrie/tslog"
```

The [global informational logger and error logger](https://pkg.go.dev/github.com/thorstenrie/tslog#pkg-variables) are used for logging.

E.g.,
```
tslog.I.Println("Hello World!") // info
tslog.E.Println("Hello Error!") // error
```

## Configuration

The tslog package can be configured to log to stdout (default), a temporary file, a specific file or logging can be discarded (no logging).

The following configurations with `TS_LOGFILE` are available

- `stdout`: Log to Stdout (default)
- `tmp`: logging to `tslog_*` in the operating system temporary directory, where `*` stands for a random string (see [os.CreateTemp](https://pkg.go.dev/os#CreateTemp))
- `discard`: no logging
- `<filename>`: for logging to <filename>

Therefore, `stdout`, `tmp`, `discard` are reserved keywords. If none of the keywords apply, the contents of `TS_LOGFILE` will be treated as filename. If `TS_LOGFILE` is not set, then tslog will fall back to the default logging to Stdout.

## Example

```
package main

import (
	"github.com/thorstenrie/tslog"
)

func main() {
	tslog.I.Println("Hello World!")
	tslog.E.Println("Hello Error!")
}
```

## Links

[Godoc](https://pkg.go.dev/github.com/thorstenrie/tslog)

[Gocover.io](https://gocover.io/github.com/thorstenrie/tslog)

[Go Report Card](https://goreportcard.com/report/github.com/thorstenrie/tslog)

[Open Source Insights](https://deps.dev/go/github.com%2Fthorstenrie%2Ftslog)
