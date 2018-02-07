FROM golang:1.9

ADD . /go/src/github.com/sharykhin/todoapp

WORKDIR /go/src/github.com/sharykhin/todoapp

RUN go get .

RUN go install github.com/sharykhin/todoapp

ENTRYPOINT /go/bin/todoapp

EXPOSE 8082