FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM debian:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/projects.json . 

EXPOSE 8080

CMD ["./main"]