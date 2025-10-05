FROM docker.io/golang:1.24.4-alpine3.22 AS builder
WORKDIR /build

RUN apk update && apk add --no-cache ca-certificates git gcc g++ libc-dev binutils
RUN --mount=type=bind,source=go.mod,target=/build/go.mod \
    --mount=type=bind,source=go.sum,target=/build/go.sum  \
    go mod download && go mod verify 

COPY ../src .
RUN go build -o bin/application .

FROM docker.io/alpine:3.21 AS runner
WORKDIR /app

RUN apk update && apk add --no-cache libc6-compat curl && rm -rf /var/cache/apk/*

COPY --from=builder /build/bin/application ./

USER guest
EXPOSE 8080
CMD ["./application"]
