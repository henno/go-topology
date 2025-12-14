# Dockerfile template - adjust for your language
#
# Go:
#   FROM golang:1.22-alpine AS builder
#   WORKDIR /app
#   COPY . .
#   RUN go build -o server ./cmd/server
#   FROM alpine
#   COPY --from=builder /app/server /server
#   CMD ["/server"]
#
# Bun:
#   FROM oven/bun:1
#   WORKDIR /app
#   COPY . .
#   RUN bun install
#   CMD ["bun", "run", "src/server.ts"]
#
# Python:
#   FROM python:3.12-slim
#   WORKDIR /app
#   COPY requirements.txt .
#   RUN pip install -r requirements.txt
#   COPY . .
#   CMD ["python", "src/server.py"]
#
# Rust:
#   FROM rust:1.75 AS builder
#   WORKDIR /app
#   COPY . .
#   RUN cargo build --release
#   FROM debian:bookworm-slim
#   COPY --from=builder /app/target/release/netmap /netmap
#   CMD ["/netmap"]

FROM alpine:3.19
RUN echo "Replace this Dockerfile with your language-specific version"
CMD ["echo", "Not configured"]
