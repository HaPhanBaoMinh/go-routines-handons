# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

# build static binary
RUN go build -o stack_test stack_probe_go.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/stack_test .

# default command (override at runtime if needed)
# Here we show an example with --ulimit in CMD so you can see it clearly
CMD ["./stack_test"]

