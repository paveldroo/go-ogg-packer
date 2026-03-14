FROM golang:1.24-trixie

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    libopus-dev=1.5.2-2 \
    libopusfile-dev=0.12-4+b3

RUN go install github.com/go-task/task/v3/cmd/task@v3.9.1

WORKDIR /app

COPY . /app

RUN go mod download
