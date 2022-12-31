FROM golang:1.19-alpine as builder
WORKDIR /root/app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/api/

FROM alpine:latest
WORKDIR /root/app
COPY --from=builder /root/app/api .
ENTRYPOINT [ "/root/app/api" ]