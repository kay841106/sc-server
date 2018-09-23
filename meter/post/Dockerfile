#!/bin/bash
FROM golang:onbuild
MAINTAINER Fransiscus Bimo <avbee.lab@gmail.com>
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o post .
CMD ["./post"]
