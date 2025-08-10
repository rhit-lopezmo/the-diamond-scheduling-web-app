# Build
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server

# Run (alpine so healthcheck can curl)
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=build /app/server /app/server
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/server"]
