# Stage 1: Build
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

# Stage 2: Run
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

# for testing purpose 
# COPY app.env .

EXPOSE 8080

CMD ["/app/server"]
