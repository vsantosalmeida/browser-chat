FROM golang:1.19-alpine3.18 as builder
WORKDIR /go/src/vsantosalmeida/browser-chat

COPY . ./

RUN apk add --no-cache \
   build-base \
   gcc \
   git \
   pkgconf \
   musl-dev

RUN export GOPRIVATE=github.com/vsantosalmeida/*
RUN make build-chat-bot

FROM alpine
WORKDIR /root/
COPY --from=builder /go/src/vsantosalmeida/browser-chat/bin/chat-bot.bin .
EXPOSE 8080
CMD ["./chat-bot.bin"]