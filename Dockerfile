FROM dhi.io/golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go

# Final chapter: THE Hardened Static Base
FROM dhi.io/static:latest

WORKDIR /app

COPY --from=builder /app/server .

# Using the pre-configured non-root user (standard for hardened images)
USER 65532

EXPOSE 8080

CMD ["./server"]
