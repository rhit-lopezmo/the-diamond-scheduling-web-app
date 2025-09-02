# Build
FROM golang:1.24-alpine AS build
WORKDIR /app

# Put the module files from api/ into the current workdir
COPY api/go.mod api/go.sum ./
RUN go mod download

# Now copy the rest of the api source into /app
COPY api/ .

# Build from /app (module root is here now)
RUN CGO_ENABLED=0 go build -o /out/server .

# Run
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache curl ca-certificates

COPY --from=build /out/server /app/server
COPY db/migrations /app/db/migrations

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/server"]
