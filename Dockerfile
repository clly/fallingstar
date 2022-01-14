FROM golang:1.17 as build

WORKDIR /build
COPY . ./
RUN make

FROM ubuntu

WORKDIR /opt
COPY --from=build /build/goregex ./

WORKDIR /data
USER nobody
ENTRYPOINT ["/opt/goregex"]
