FROM balenalib/%%BALENA_MACHINE_NAME%%-golang:latest-build AS build

WORKDIR /go/src/github.com/stephenhillier/geoprojects/devices/therm

COPY /app ./

RUN go get ./... && go build

FROM balenalib/%%BALENA_MACHINE_NAME%%-debian:stretch

COPY --from=build /go/src/github.com/stephenhillier/geoprojects/devices/therm/therm .

CMD ./therm