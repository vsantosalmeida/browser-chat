FROM golang:1.19-alpine3.18 as builder
WORKDIR /app

COPY . .

RUN apk add --no-cache \
   build-base \
   gcc \
   git \
   pkgconf \
   musl-dev

RUN make build-chat-api

FROM alpine
WORKDIR /app
COPY --from=builder /app/bin/chat-api.bin .
COPY --from=builder /app/templates templates/
EXPOSE 8080
EXPOSE 3000
CMD ["./chat-api.bin"]