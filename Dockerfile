FROM golang:1.21-alpine3.17 AS builder

ARG ZIG_RELEASE="0.10.1"

WORKDIR /src

# install git to inject version into binary
RUN apk add -U --no-cache \
    git

# install zig to build sqlite3 with cgo
# for some reason, zig 0.11.0 can't compile sqlite3, I had to downgrade to 0.10.1
RUN wget -qO- https://ziglang.org/download/${ZIG_RELEASE}/zig-linux-x86_64-${ZIG_RELEASE}.tar.xz | tar -xJ -C /opt

ENV CC="/opt/zig-linux-x86_64-${ZIG_RELEASE}/zig cc -target native-native-musl"

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 go build -o /scwtui ./cmd/scwtui

FROM alpine:3.17

# this is needed for the colors to display correctly
ENV TERM="xterm-256color"

RUN apk add -U --no-cache ca-certificates

COPY --from=builder /scwtui /scwtui

ENTRYPOINT ["/scwtui"]
