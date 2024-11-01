FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o ticket-booking

FROM alpine:latest

COPY --from=builder /app/ticket-booking /ticket-booking

EXPOSE 8080

CMD ["/ticket-booking"]