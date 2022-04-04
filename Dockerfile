## Stage 1
FROM golang:1.17-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o speedtest-exporter

## Stage 2
FROM debian:11.1-slim

ARG SPEEDTEST_VERSION=speedtest_1.0.0.2-1.5ae238b

RUN apt-get update && \
    apt-get install -y wget && \
    rm -rf /var/lib/apt/lists/* && \
    export ARCH=$(dpkg --print-architecture) && \
    wget --content-disposition https://packagecloud.io/ookla/speedtest-cli/packages/debian/bullseye/${SPEEDTEST_VERSION}_${ARCH}.deb/download.deb && \
    dpkg -i ./${SPEEDTEST_VERSION}_${ARCH}.deb && \
    rm -f ${SPEEDTEST_VERSION}_${ARCH}.deb

# Accept license, speedtest will automaticly start the test, 
# in case that the program fails, ignore error by || true

# RUN speedtest --accept-license --accept-gdpr  || true

COPY --from=builder /app/speedtest-exporter /app/speedtest-exporter
CMD ["/app/speedtest-exporter"]
