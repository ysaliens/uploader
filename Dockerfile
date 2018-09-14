FROM golang:1.11 as builder
COPY ./handlers /go/src/github.com/ysaliens/uploader/handlers
COPY ./models /go/src/github.com/ysaliens/uploader/models
ARG SOURCE_LOCATION=/
WORKDIR ${SOURCE_LOCATION}
RUN go get -d -v github.com/gin-gonic/gin \
    && go get -d -v github.com/tealeg/xlsx \
	&& go get -d -v gopkg.in/mgo.v2/bson \
	&& go get -d -v gopkg.in/mgo.v2
COPY uploader.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
ARG SOURCE_LOCATION=/
RUN apk --no-cache add curl
EXPOSE 8080
WORKDIR /root/
COPY --from=builder ${SOURCE_LOCATION} .
CMD ["./app"]  