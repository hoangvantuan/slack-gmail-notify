GPATH=$(GOPATH)

tunnel:
	tunnel -config $(GPATH)/src/github.com/mdshun/slack-gmail-notify/.tunnel/tunnel.yml start-all
dev:
	gin -x vendor -i --appPort 8081 --port 8080 --path $(GPATH)/src/github.com/mdshun/slack-gmail-notify run main.go
stg:
	gin -x vendor -i --appPort 8082 --port 8083 --path $(GPATH)/src/github.com/mdshun/slack-gmail-notify run main.go
run:
	go install
	./slack-gmail-notify > dev/null 2>&1 &
