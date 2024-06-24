FROM golang AS builder
WORKDIR /source
RUN apt-get update && apt-get install -y git
RUN git clone https://github.com/Weebs-Kingdom/weebskingdom
WORKDIR ./weebskingdom/main
RUN go mod download
RUN go build -o /app/weebskingdom

FROM frolvlad/alpine-glibc

WORKDIR /app

COPY --from=builder "/app/weebskingdom" .

RUN mkdir /config

VOLUME /config

EXPOSE 80

ENTRYPOINT ["/app/weebskingdom"]
