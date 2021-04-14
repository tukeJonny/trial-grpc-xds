FROM golang:latest as builder

WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build -o server .

FROM debian:buster
COPY --from=builder /build/server .

EXPOSE 8080
ENTRYPOINT ["./server"]