package main

import (
	"flag"
	"fmt"

	"github.com/m4gshm/flag/flagenum"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

func main() {
	var (
		toEnum      = func(s string) Enum { return Enum(Enum_value[s]) }
		toStr       = Enum.String
		values, err = flagenum.Multiple(
			flagenum.CommandLine.FlagSet,
			"enum",
			slice.Of(Enum_A, Enum_D), /*default*/
			slice.Convert(map_.Keys(Enum_value), toEnum), /*allowed*/
			toEnum, toStr, "grpc enum example",
		)
	)
	slice.Convert(map_.Keys(Enum_value), toEnum)

	map_.ToSlice(Enum_value, func(k string, v int32) Enum { return toEnum(k) })

	if err != nil {
		panic(err)
	}
	flag.Parse()

	fmt.Printf("enum values: %v\n", *values)
}
