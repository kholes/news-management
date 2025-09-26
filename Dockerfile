FROM golang:1.24-alpine AS builder

WORKDIR /app


# Install git for go get (modules)
RUN apk add gcc musl-dev --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /news-api ./main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /news-api /news-api
COPY .env /app/.env
EXPOSE 8080
ENTRYPOINT ["/news-api"]