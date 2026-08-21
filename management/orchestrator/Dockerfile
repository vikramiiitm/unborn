FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /orchestrator ./main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates docker-cli android-tools
WORKDIR /app
COPY --from=builder /orchestrator /app/orchestrator
EXPOSE 8080
ENTRYPOINT ["/app/orchestrator"]
