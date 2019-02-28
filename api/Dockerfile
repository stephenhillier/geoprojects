FROM golang:1.11-stretch AS builder
WORKDIR /go/src/github.com/stephenhillier/geoprojects/api
ADD . /go/src/github.com/stephenhillier/geoprojects/api/
RUN go test && CGO_ENABLED=0 go build

FROM alpine:3.9
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /go/src/github.com/stephenhillier/geoprojects/api/api /app/
ENTRYPOINT [ "/app/api" ]
EXPOSE 8000
