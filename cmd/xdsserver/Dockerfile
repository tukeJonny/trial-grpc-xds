FROM golang:latest as builder

WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build -o xdsserver .

FROM debian:buster
COPY --from=builder /build/xdsserver .

EXPOSE 8080
ENTRYPOINT ["./xdsserver"]