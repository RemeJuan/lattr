FROM golang:1.16.6-alpine AS builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build -o lattr main.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/lattr /app/
WORKDIR /app
CMD ["./lattr"]
