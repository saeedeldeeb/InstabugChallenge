FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0  \
    GOARCH="amd64" \
    GOOS=linux

# Install git
RUN apk add --no-cache git

WORKDIR /build
COPY . .

RUN go mod tidy

# Only build the application
RUN go build --ldflags "-extldflags -static" -o main .

FROM alpine:latest

WORKDIR /www

COPY --from=builder /build/main /www/
COPY --from=builder /build/database/ /www/database/
COPY --from=builder /build/public/ /www/public/
COPY --from=builder /build/storage/ /www/storage/
COPY --from=builder /build/resources/ /www/resources/
COPY --from=builder /build/.env.example /www/.env

# Create an entrypoint script
COPY docker-entrypoint.sh /www/
RUN chmod +x /www/docker-entrypoint.sh

ENTRYPOINT ["/www/docker-entrypoint.sh"]
