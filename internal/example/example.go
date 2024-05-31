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
