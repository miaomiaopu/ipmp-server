FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILDTIME=unknown
RUN CGO_ENABLED=0 go build \
  -ldflags "-X github.com/miaomiaopu/ipmp-server/internal/pkg/version.Version=${VERSION} -X github.com/miaomiaopu/ipmp-server/internal/pkg/version.GitCommit=${COMMIT} -X github.com/miaomiaopu/ipmp-server/internal/pkg/version.BuildTime=${BUILDTIME}" \
  -o ipmp-server cmd/server/main.go

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/ipmp-server /usr/local/bin/
EXPOSE 8080
CMD ["ipmp-server"]
