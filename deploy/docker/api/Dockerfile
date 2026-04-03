FROM golang:latest AS build

WORKDIR /app

RUN apt-get update \
    && apt-get install -y jq \
    && rm -rf /var/lib/apt/lists/*

# Run download separately to leverage docker cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./build/api ./cmd/api

RUN jq -n --arg commit "$(git rev-parse --short HEAD)" \
        --arg timestamp "$(date '+%Y%m%dT%H%M%S')" \
        '{commit: $commit, timestamp: $timestamp}' > "deploy.json"


FROM ubuntu:24.04 AS runtime

WORKDIR /app

RUN apt-get update \
    && apt-get install -y tzdata ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/build/api /app/api
COPY --from=build /app/deploy.json /app/deploy.json

RUN chmod +x api

CMD ["./api"]
