#!/bin/bash

# Fail explicitly if needed
set -euxo pipefail

# Clone ffmpeg
if [ ! -d "./ffmpeg" ]; then
      git clone --depth 1 git://source.ffmpeg.org/ffmpeg
fi

# Build ffmpeg
cd ffmpeg
PKG_CONFIG_PATH="/usr/local/lib/pkgconfig" \
      ./configure --extra-libs=-lpthread --prefix="/usr/local" \
      --extra-cflags="-I/usr/local/include" --extra-ldflags="-L/usr/local/lib" \
      --pkg-config-flags="--static" --enable-gpl --enable-nonfree \
      --enable-libfdk-aac --enable-libmp3lame --enable-libopus \
      --enable-libvorbis --enable-libvpx --enable-libx264 --enable-libx265 \
      --enable-libaom
make -j $(nproc)

# Install ffmpeg
make install
hash -r