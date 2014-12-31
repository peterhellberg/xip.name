all: build

build:
	go build -o xip xip.go

linux:
	GOOS=linux CGO_ENABLED=0 go build -o xip.linux xip.go

run:
	go run xip.go -p 8053 -v

dig:
	dig @localhost -p 8053 foo.bar.10.1.2.3.xip.name A

clean:
	rm -f xip xip.linux
