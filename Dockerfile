FROM golang:1.21.5-alpine AS builder
WORKDIR /app
COPY cmd/ ./cmd/
COPY internal/ ./internal
COPY proto/ ./proto
COPY go.mod .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/server/

FROM alpine AS runner
WORKDIR /
COPY --from=builder /app/main /main
COPY --from=builder /app/internal/template/html /html

ENTRYPOINT ["./main"]