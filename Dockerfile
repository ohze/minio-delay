FROM alpine:3.6

MAINTAINER Ohze JSC <thanhbv@sandinh.net>

ENV GOPATH=/go \
 PATH=$PATH:/go/bin \
 CGO_ENABLED=0 \
 MINIO_ENDPOINT=lb.minio:9000/ \
 MINIO_KEY=minio \
 MINIO_SECRET=minio123


WORKDIR /go/src/github.com/ohze/

RUN \
    apk add --no-cache ca-certificates && \
    apk add --no-cache --virtual .build-deps go git musl-dev && \
    echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

    go get -v -d github.com/ohze/minioxf && \
    cd /go/src/github.com/ohze/minioxf && \
    go install -v  -ldflags "$(go run buildscripts/gen-ldflags.go)" && \
    rm -rf /go/pkg /go/src /usr/local/go && apk del .build-deps


EXPOSE 9009

CMD ["minioxf"]
