FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o runlog-api ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/runlog-api .

EXPOSE 8080

CMD ["./runlog-api"]