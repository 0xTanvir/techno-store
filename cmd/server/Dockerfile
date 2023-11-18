# Build stage
FROM golang:1.21-alpine3.18 AS build

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/server/ -a -installsuffix cgo -o server .

# Final stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /app/cmd/server .
COPY app.env .

EXPOSE 8080
CMD ["./server"]
