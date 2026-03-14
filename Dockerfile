FROM golang:1.19-bullseye

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    libopus-dev=1.3.1-0.1 \
    libopusfile-dev=0.9+20170913-1.1 \
    ffmpeg

RUN go install github.com/go-task/task/v3/cmd/task@v3.9.1

WORKDIR /app

COPY . /app

RUN go mod download
