FROM golang:1.19-alpine3.18 as builder
WORKDIR /app

COPY . .

RUN apk add --no-cache \
   build-base \
   gcc \
   git \
   pkgconf \
   musl-dev

RUN make build-chat-bot

FROM alpine
WORKDIR /root/
COPY --from=builder /app/bin/chat-bot.bin .
CMD ["./chat-bot.bin"]