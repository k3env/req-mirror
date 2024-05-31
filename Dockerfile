FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o mirror main.go

FROM alpine
WORKDIR /app
COPY --from=builder /build/mirror /app/mirror
ENTRYPOINT ["/app/mirror"]
