
FROM scratch
MAINTAINER Fransiscus Bimo <avbee.lab@gmail.com>
ENV MONGO_URL "172.16.0.132:27017"
ENV USERNAME "dontask"
ENV PASS "idontknow"
ENV DB_AUTH "admin"
ADD post /
ENTRYPOINT ["/post"]


#FROM golang:onbuild
#MAINTAINER Fransiscus Bimo <avbee.lab@gmail.com>
#RUN mkdir /app
#ADD . /app/ 
#WORKDIR /app
#RUN go build -o post .
#CMD ["/app/post"]

