module example

go 1.22

toolchain go1.22.1

replace github.com/m4gshm/flag v0.0.0-20240529065328-b85af3c22117 => ../../

require github.com/m4gshm/flag v0.0.0-20240529065328-b85af3c22117

require (
	github.com/m4gshm/gollections v0.0.12
	golang.org/x/exp v0.0.0-20240531132922-fd00a4e0eefc
	google.golang.org/protobuf v1.34.1 // indirect
)
