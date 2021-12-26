# Build 
FROM golang:1.15.13 AS build
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.* /app/
RUN go mod download

COPY . .
RUN go build -o /payment-service

# Deploy
FROM alpine:3.14

WORKDIR /

COPY --from=build /payment-service /payment-service
COPY .env .env

CMD ["./payment-service"]
