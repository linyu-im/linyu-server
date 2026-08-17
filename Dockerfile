FROM golang:1.26-alpine AS builder

WORKDIR /src

ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn

RUN apk add --no-cache git ca-certificates tzdata

COPY go.work go.work.sum ./
COPY go.mod go.sum ./
COPY linyu-ai/go.mod linyu-ai/go.sum ./linyu-ai/
COPY linyu-application/go.mod linyu-application/go.sum ./linyu-application/
COPY linyu-auth/go.mod linyu-auth/go.sum ./linyu-auth/
COPY linyu-basic-service/go.mod linyu-basic-service/go.sum ./linyu-basic-service/
COPY linyu-cloud-drive/go.mod linyu-cloud-drive/go.sum ./linyu-cloud-drive/
COPY linyu-common/go.mod linyu-common/go.sum ./linyu-common/
COPY linyu-gateway/go.mod linyu-gateway/go.sum ./linyu-gateway/
COPY linyu-im/go.mod linyu-im/go.sum ./linyu-im/
COPY linyu-voip-chat/go.mod linyu-voip-chat/go.sum ./linyu-voip-chat/

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY mian.go ./
COPY assets ./assets
COPY linyu-ai ./linyu-ai
COPY linyu-application ./linyu-application
COPY linyu-auth ./linyu-auth
COPY linyu-basic-service ./linyu-basic-service
COPY linyu-cloud-drive ./linyu-cloud-drive
COPY linyu-common ./linyu-common
COPY linyu-gateway ./linyu-gateway
COPY linyu-im ./linyu-im
COPY linyu-voip-chat ./linyu-voip-chat

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/linyu-server .

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 10001 linyu

COPY --from=builder /out/linyu-server /app/linyu-server
COPY --from=builder /src/assets /app/assets

RUN chown -R linyu:linyu /app

USER linyu

EXPOSE 9088

VOLUME ["/app/assets/config", "/app/linyu"]

ENTRYPOINT ["/app/linyu-server"]
CMD ["-config", "assets/config/config.yml", "-locales", "assets/locales", "-email-templates", "assets/email_templates"]
