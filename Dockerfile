FROM golang

ARG env=dev
ARG app=${GOPATH}/src/github.com/mdshun/slack-gmail-notify
ARG port=8080
ENV GOBIN=$GOPATH/bin

RUN go get github.com/mdshun/slack-gmail-notify
RUN go get github.com/golang/dep/cmd/dep
WORKDIR ${app}
RUN go install
EXPOSE ${port}

#CMD ["go", "run", "main.go"]
