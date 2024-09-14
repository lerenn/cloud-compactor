# Building ffmpeg with extensions
################################################################################
FROM alpine:latest AS ffmpeg-build

# Install requirements
RUN apk add \
    autoconf \
    automake \
    bash \
    cmake \
    curl \
    freetype-dev \
    numactl-dev \
    g++ \
    gcc \
    gettext-dev \
    git \
    libogg-dev \
    libtool \
    libvorbis-dev \
    libvpx-dev \
    make \
    mercurial \
    nasm \
    opus-dev \
    pkgconfig \
    yasm \
    zlib-dev

# Set workspace
RUN mkdir -p /build /scripts
WORKDIR /build 

# Compile x264
COPY ./build/package/install/x264.sh /scripts/x264.sh
RUN bash /scripts/x264.sh

# Compile x265
COPY ./build/package/install/x265.sh /scripts/x265.sh
RUN bash /scripts/x265.sh

# Compile libaom
COPY ./build/package/install/libaom-v1.sh /scripts/libaom-v1.sh
RUN bash /scripts/libaom-v1.sh

# Compile libfdk_aac
COPY ./build/package/install/libfdk_aac.sh /scripts/libfdk_aac.sh
RUN bash /scripts/libfdk_aac.sh

# Compile libmp3lame
COPY ./build/package/install/libmp3lame.sh /scripts/libmp3lame.sh
RUN bash /scripts/libmp3lame.sh

# Compile ffmpeg
COPY ./build/package/install/ffmpeg.sh /scripts/ffmpeg.sh
RUN bash /scripts/ffmpeg.sh

# Building go binary
################################################################################
FROM golang:alpine AS go-build

# Disable CGO
ENV CGO_ENABLED=0

# Get all remaining code
RUN mkdir -p /go/src/github.com/lerenn/cloud-compactor
COPY ./ /go/src/github.com/lerenn/cloud-compactor

# Set the workdir
WORKDIR /go/src/github.com/lerenn/cloud-compactor

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go install ./cmd/*

# Final base image
################################################################################
FROM alpine:latest

# Install requirements
RUN apk add libvpx opus-dev libvorbis libstdc++ numactl-dev

# Set environment variables
ENV HEALTH_PORT=9000

# Expose ports
EXPOSE 9000

# Get binary
COPY --from=ffmpeg-build /usr/local/bin/* /usr/local/bin
COPY --from=ffmpeg-build /usr/local/lib/* /usr/local/lib
COPY --from=go-build     /go/bin/*        /usr/local/bin

# Set entrypoint
ENTRYPOINT [ "cloud-compactor" ]