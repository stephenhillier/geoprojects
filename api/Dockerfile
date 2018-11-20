FROM golang:1.11-stretch AS builder
ENV GO111MODULE on
WORKDIR /go/src/github.com/stephenhillier/geoprojects/api
ADD go.* /go/src/github.com/stephenhillier/geoprojects/api/
RUN go mod download
ADD . /go/src/github.com/stephenhillier/geoprojects/api/
RUN go test && go build

FROM centos
WORKDIR /app
COPY --from=builder /go/src/github.com/stephenhillier/geoprojects/api/api /app/
RUN chmod +x /app/api
ENTRYPOINT [ "/app/api" ]
EXPOSE 8000
