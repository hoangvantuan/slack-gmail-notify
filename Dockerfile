FROM golang

ENV SLGMAILS_ENV=dev
ENV GOBIN=$GOPATH/bin
ENV app=$GOPATH/src/github.com/mdshun/slack-gmail-notify

WORKDIR ${app}

RUN export PATH=$PATH:$GOBIN && \
    git clone https://github.com/mdshun/slack-gmail-notify ${app} && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    go build -o sgn main.go

COPY *.env ./
EXPOSE 8080

CMD [ "./sgn" ]