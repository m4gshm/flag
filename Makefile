.PHONY: all
all: build test readme lint

.PHONY: test
test:
	$(info #Running tests...)
	go clean -testcache
	go test ./...

.PHONY: build
build:
	$(info #Building...)
	go build ./...

.PHONY: lint
lint:
	$(info #Lints...)
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	go vet ./...
	go install github.com/tetafro/godot/cmd/godot@latest
	godot .
	go install github.com/kisielk/errcheck@latest
	errcheck ./...
	go install golang.org/x/lint/golint@latest
	golint ./...

.PHONY: readme
readme:
	$(info #README.md...)
	cd internal/example && go run . > ../docs/run1.txt
	cd internal/example && go run . --api soap --api rest --log-level debug > ../docs/run2.txt
	cd internal/example && go run . --help > ../docs/usage.txt 2>&1 && tail -n +2 "../docs/usage.txt" > "../docs/usage.tmp" && mv ../docs/usage.tmp ../docs/usage.txt
	asciidoctor -b docbook internal/docs/readme.adoc 
	pandoc -f docbook -t gfm internal/docs/readme.xml -o README.md