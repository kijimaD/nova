########
# base #
########

# なぜかbuster以外だと、WASMビルドで真っ白表示になってしまう
FROM golang:1.22.8-bullseye AS base
RUN apt update
RUN apt install -y \
    gcc \
    libc6-dev \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxxf86vm-dev \
    libasound2-dev \
    pkg-config \
    xorg-dev \
    libx11-dev \
    libopenal-dev \
    upx-ucl

###########
# builder #
###########

FROM base AS builder

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GO111MODULE=on \
    go build -o ./bin/nova ./_example/main.go
RUN upx-ucl --best --ultra-brute ./bin/nova

###########
# release #
###########

FROM gcr.io/distroless/base-debian11:latest AS release

COPY --from=builder /build/bin/nova /bin/
WORKDIR /work
ENTRYPOINT ["nova"]

########
# node #
########

FROM node:22 as releaser
RUN yarn install
