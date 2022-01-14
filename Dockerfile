FROM golang:1.17 as build

WORKDIR /build
COPY . ./
RUN make

FROM ubuntu

WORKDIR /opt
COPY --from=build /build/fallingstar ./

RUN apt-get update && apt-get install -y \
  ca-certificates \
  git \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /data
USER nobody
ENTRYPOINT ["/opt/fallingstar"]
