FROM golang:1.21 as builder

WORKDIR /app/

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o application cmd/gateway/main.go

FROM alpine:3.18.6

WORKDIR /app/
COPY --from=builder /app/application ./
COPY --from=builder /app/configs/gateway/*.yaml ./

CMD ["/app/application", "-c", "/app/config.yaml"]
