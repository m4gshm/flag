# Flag

Extends the core flag package by adding implementations, such as an
argument with multiple values or an argument limited to predefined
options.

## Example

source `example.go`:

``` go
package main

import (
    "flag"
    "fmt"

    "github.com/m4gshm/flag/flagenum"
)

func main() {
    var (
        api = flagenum.Strings(
            "api",
            []string{"rest", "grpc"},         /*default*/
            []string{"rest", "grpc", "soap"}, /*allowed*/
            "enabled api engine",
        )
        logLevel = flagenum.String(
            "log-level",
            "info", /*default*/
            []string{"debug", "info", "warn", "error"}, /*allowed*/
            "logger level",
        )
    )
    flag.Parse()

    fmt.Printf("enabled apis: %v\n", *api)
    fmt.Printf("log level:    %s\n", *logLevel)
}
```

Command `go run .` will print:

``` console
enabled apis: [rest grpc]
log level:    info
```

To change defaults need set arguments like
`go run . --api soap --api rest --log-level debug`

``` console
enabled apis: [soap rest]
log level:    debug
```

Call `go run . --help` to get usage info:

``` console
Usage of example:
  -api value
        enabled api engine (allowed any of rest,grpc,soap) (default rest,grpc)
  -log-level value
        logger level (allowed one of debug,info,warn,error) (default info)
```
