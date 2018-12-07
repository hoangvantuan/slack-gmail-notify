GPATH=$(GOPATH)

tunnel:
	tunnel -config $(GPATH)/src/github.com/mdshun/slack-gmail-notify/.tunnel/tunnel.yml start-all
dev:
	gin -x vendor --appPort 8081 --port 8080 --path $(GPATH)/src/github.com/mdshun/slack-gmail-notify run main.go
