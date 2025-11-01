FROM golang:1.22-alpine AS build

RUN apk add --no-cache build-base git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

FROM alpine:latest AS runtime
RUN addgroup -S app && adduser -S app -G app

WORKDIR /srv/app

COPY --from=build /app/server ./server
COPY configs ./configs

USER app
EXPOSE 8080

ENTRYPOINT ["./server"]
