FROM golang:1.23.2-bookworm AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod verify
RUN go mod download

COPY . .

RUN go build .

FROM debian:bookworm-slim

COPY --from=builder /app/backend .

CMD ["./backend"]
