FROM golang AS builder
WORKDIR /source
RUN apt-get update && apt-get install -y curl
ADD "https://api.github.com/repos/Weebs-Kingdom/weebskingdom/commits?per_page=1" latest_commit
RUN curl -sL "https://github.com/Weebs-Kingdom/weebskingdom/archive/main.zip" -o weebskingdom.zip && unzip weebskingdom.zip
WORKDIR ./weebskingdom-main/main
RUN go mod download
RUN go build -o /app/weebskingdom

FROM frolvlad/alpine-glibc

WORKDIR /app

COPY --from=builder "/app/weebskingdom" .

RUN mkdir /config

VOLUME /config

EXPOSE 80

ENTRYPOINT ["/app/weebskingdom"]
