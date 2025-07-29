FROM golang:1.20

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    libopus-dev=1.3.1-3 \
    libopusfile-dev=0.12-4

RUN go install github.com/go-task/task/v3/cmd/task@v3.9.1

WORKDIR /app

COPY . /app

RUN go mod download
