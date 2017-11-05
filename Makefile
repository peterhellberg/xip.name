all: build

build:
	go build -o xip xip.go

linux:
	GOOS=linux CGO_ENABLED=0 go build -o xip.linux xip.go

init:
	scp etc/init/xip.name.conf root@188.166.43.179:/etc/init/

web:
	scp usr/share/nginx/html/* root@188.166.43.179:/usr/share/nginx/html/

deploy: linux init web
	ssh root@188.166.43.179 'service xip.name stop || true'
	scp xip.linux root@188.166.43.179:/usr/local/bin/xip.name
	ssh root@188.166.43.179 'service xip.name start'

run:
	go run xip.go -p 8053 -v

dig:
	dig @localhost -p 8053 foo.bar.10.1.2.3.xip.name A

clean:
	rm -f xip xip.linux
