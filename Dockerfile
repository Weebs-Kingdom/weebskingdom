FROM golang AS builder
WORKDIR /source
RUN apt-get update && apt-get install -y git
RUN git clone https://github.com/Weebs-Kingdom/weebskingdom
WORKDIR ./weebkingdom/main
RUN go mod download
RUN go build -o /app/weesbskingdom

FROM frolvlad/alpine-glibc

WORKDIR /app

COPY --from=builder "/app/weebskingdom" .

RUN mkdir /config

VOLUME /config

EXPOSE 80

ENTRYPOINT ["/app/weebskingdom"]
