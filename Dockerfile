FROM golang:1.24

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    libopus-dev \
    libopusfile-dev

RUN go install github.com/go-task/task/v3/cmd/task@v3.9.1

WORKDIR /app

COPY . /app

RUN go mod download
