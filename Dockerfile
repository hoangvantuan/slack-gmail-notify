FROM golang

ENV SLGMAILS_ENV=dev
ENV GOBIN=$GOPATH/bin
ENV app=$GOPATH/src/github.com/mdshun/slack-gmail-notify

RUN export PATH=$PATH:$GOBIN
RUN git clone https://github.com/mdshun/slack-gmail-notify ${app}
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${app}
COPY *.env ./
RUN go build -o sgn main.go
EXPOSE 8080
CMD [ "./sgn" ]