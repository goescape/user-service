FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY config.yaml.docker config.yaml

RUN echo "Building the application..."

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./main.go

# Build a small image

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

# CMD ["/app/server"]
CMD [ "sh", "-c", "sleep 20 && /app/server" ]
