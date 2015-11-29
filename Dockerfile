FROM golang:1.5

COPY . /opt/pault.ag/deceive
WORKDIR /opt/pault.ag/deceive

RUN go get -d .
RUN go build -o /usr/local/bin/deceive .
