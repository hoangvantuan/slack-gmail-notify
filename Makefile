GPATH=$(GOPATH)

dev:
	gin -x vendor -i --appPort 8081 --port 8080 --path $(GPATH)/src/github.com/mdshun/slack-gmail-notify run main.go
stg:
	kill -KILL $(cat ./stg.pid)
	go get -u github.com/mdshun/slack-gmail-notify
	SLGMAIL_ENV=stg slack-gmail-notify > /dev/null 2>&1 &
	echo $! > ./stg.pid
create-tunnel-client:
	openssl req -x509 -nodes -newkey rsa:2048 -sha256 -keyout .tunnel/client.key -out .tunnel/client.crt
create-tunnel-server:
	openssl req -x509 -nodes -newkey rsa:2048 -sha256 -keyout .tunnel/server.key -out .tunnel/server.crt
tunnel-client:
	tunnel -config .tunnel/tunnel.yml start-all
tunnel-server:
	kill -KILL $(cat ./tunnel-server.pid)
	tunneld -httpAddr :5000 -httpsAddr :5001 -tlsCrt .tunnel/server.crt -tlsKey .tunnel/server.key > /dev/null 2>&1 &
	echo $! > ./tunnel-server.pid
.PHONY: caddy
caddy:
	kill -KILL $(cat ./caddy.pid)
	caddy > /dev/null 2>&1 &
	echo $! > ./caddy.pid