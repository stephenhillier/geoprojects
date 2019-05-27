FROM golang:1.12-stretch AS builder
WORKDIR /go/src/github.com/stephenhillier/geoprojects/earthworks
ADD . /go/src/github.com/stephenhillier/geoprojects/earthworks/

# run tests and, if they pass, build the server
RUN go get ./... && go test ./... && cd ./cmd/server && CGO_ENABLED=0 go build
RUN CGO_ENABLED=0 go get -u github.com/amacneil/dbmate

FROM alpine:3.9
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /go/bin/dbmate /app/
COPY --from=builder /go/src/github.com/stephenhillier/geoprojects/earthworks/cmd/server /app/
COPY --from=builder /go/src/github.com/stephenhillier/geoprojects/earthworks/db/migrations/* /app/db/migrations/
ENTRYPOINT [ "/app/server" ]
EXPOSE 8000
