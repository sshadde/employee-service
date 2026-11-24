FROM golang:1.27-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

FROM alpine:3.22
RUN apk --no-cache add ca-certificates && adduser -D -H appuser
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080 9090
USER appuser
ENTRYPOINT ["./server"]