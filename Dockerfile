FROM golang:1.23.2-bookworm

RUN --mount=type=cache,target=/var/cache/apt \
  apt-get update && apt-get install -y build-essential

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod verify
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]
