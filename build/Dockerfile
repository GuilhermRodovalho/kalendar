FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /kalendar-api ./cmd/api

FROM scratch

COPY --from=builder /kalendar-api /kalendar-api

EXPOSE 8080

ENTRYPOINT ["/kalendar-api"]
