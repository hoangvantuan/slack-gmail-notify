FROM golang

ENV GOBIN=$GOPATH/bin
ENV app=/user/app

WORKDIR ${app}

RUN export PATH=$PATH:$GOBIN && \
    go get github.com/mdshun/slack-gmail-notify

COPY *.env ./
EXPOSE 8081

CMD [ "slack-gmail-notify" ]