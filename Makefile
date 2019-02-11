.PHONY: test

WORKDIR=`pwd`

default: build

install:
	go get github.com/smallnest/gen

vet:
	go vet .

tools:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u honnef.co/go/tools/cmd/gosimple
	go get -u honnef.co/go/tools/cmd/unused
	go get -u github.com/gordonklaus/ineffassign
	go get -u github.com/fzipp/gocyclo
	go get -u github.com/golang/lint/golint

lint:
	golint ./...

staticcheck:
	staticcheck -ignore "$(shell cat .checkignore)" .

gosimple:
	# gosimple -ignore "$(shell cat .gosimpleignore)" .
	gosimple .

unused:
	unused .

gocyclo:
	@ gocyclo -over 20 $(shell find . -name "*.go" |egrep -v "pb\.go|_test\.go")

check: staticcheck gosimple unused gocyclo

doc:
	godoc -http=:6060

deps:
	go list -f '{{ join .Deps  "\n"}}' . |grep "/" | grep -v "github.com/smallnest/gen"| grep "\." | sort |uniq

fmt:
	go fmt .

build:
	go build .

test:
	go test  -v ./test
run:
	# go run main.go --connstr "root:123456@tcp(127.0.0.1:3306)/test?&parseTime=True" --database test  --json --gorm 
	@echo "generating..."
	@go run main.go 
	@echo "generated."
