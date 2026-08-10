# syntax=docker/dockerfile:1.7

FROM node:24-alpine AS web-builder
WORKDIR /src

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY jsconfig.json svelte.config.js vite.config.js tokens.css ./
COPY src ./src
COPY static ./static
RUN npm run build:web

FROM golang:1.26-alpine AS go-builder
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY *.go ./
COPY internal ./internal
COPY --from=web-builder /src/dist ./dist
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /badge-api .

FROM scratch
WORKDIR /app

COPY --from=go-builder /badge-api ./badge-api

USER 65532:65532
EXPOSE 8080
ENV BADGE_API_HTTP_ADDR=:8080

ENTRYPOINT ["/app/badge-api"]
