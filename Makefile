all: build

build:
	go build -o xip xip.go

linux:
	GOOS=linux CGO_ENABLED=0 go build -o xip.linux xip.go

init:
	scp etc/init/xip.name.conf root@xip.name:/etc/init/

web:
	scp usr/share/nginx/html/* root@xip.name:/usr/share/nginx/html/

deploy: linux init web
	ssh root@xip.name 'service xip.name stop || true'
	scp xip.linux root@xip.name:/usr/local/bin/xip.name
	ssh root@xip.name 'service xip.name start'

run:
	go run xip.go -p 8053 -v

dig:
	dig @localhost -p 8053 foo.bar.10.1.2.3.xip.name A

clean:
	rm -f xip xip.linux
