FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . .

ARG TARGET_CMD

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/${TARGET_CMD}

FROM gcr.io/distroless/base

COPY --from=builder /out/app /app

ENTRYPOINT ["/app"]