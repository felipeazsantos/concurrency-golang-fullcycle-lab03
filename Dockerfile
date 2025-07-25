FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/auction
RUN go build -o auction

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/cmd/auction/auction ./auction
COPY cmd/auction/.env .env
EXPOSE 8080
CMD ["./auction"]
