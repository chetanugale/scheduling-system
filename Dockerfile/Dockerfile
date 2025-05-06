FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o scheduling-system

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/scheduling-system .
EXPOSE 8080
CMD ["./scheduling-system"]
