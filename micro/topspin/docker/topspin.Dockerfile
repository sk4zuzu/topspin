FROM golang:1.10-alpine3.8 as build

RUN apk --no-cache add protobuf

RUN apk --no-cache add git \
 && go get github.com/Masterminds/glide

ENV SOURCE ${GOPATH}/src/github.com/sk4zuzu/topspin/micro/topspin

COPY glide.* ${SOURCE}/

RUN cd ${SOURCE}/ && glide install \
 && cd ${SOURCE}/vendor/github.com/micro/protobuf/protoc-gen-go/ && go build \
 && mv ${SOURCE}/vendor/github.com/micro/protobuf/protoc-gen-go/protoc-gen-go ${GOPATH}/bin

COPY cmd/ ${SOURCE}/cmd/
COPY proto/ ${SOURCE}/proto/
COPY api/ ${SOURCE}/api/
COPY srv1/ ${SOURCE}/srv1/
COPY srv2/ ${SOURCE}/srv2/

RUN cd ${SOURCE}/cmd/topspin/ \
 && go generate \
 && go build

RUN cd ${SOURCE}/ \
 && go test ./...

RUN mv ${SOURCE}/cmd/topspin/topspin /usr/local/bin/

FROM alpine:3.8

RUN apk --no-cache add curl

COPY --from=build /usr/local/bin/topspin /usr/local/bin/topspin

CMD /usr/local/bin/topspin
