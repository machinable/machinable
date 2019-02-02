# golang:alpine is the alpine image with the go tools added.. manually add git
FROM golang:alpine as builder

# install gcc for compilation
RUN apk add --update gcc musl-dev
# Set an env var that matches github repo name
# ENV CGO_ENABLED=0
ENV SRC_DIR=${HOME}/go/src/github.com/anothrnick/machinable/

# Add the source code:
ADD . $SRC_DIR

# Build it:
# NOTE - dependencies will need to be in the vendor directory before building this image
RUN cd $SRC_DIR;\
    apk add --no-cache git;\
    go build -o api;

ENTRYPOINT ["/go/src/github.com/anothrnick/machinable/api"]

# alpine production environment
# copy binary for smallest image size
FROM alpine:3.7

RUN apk add --no-cache ca-certificates

ENV GIN_MODE=release

COPY --from=builder /go/src/github.com/anothrnick/machinable/api /bin/api

ENTRYPOINT ["/bin/api"]
