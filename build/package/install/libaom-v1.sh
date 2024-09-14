#!/bin/bash

# Fail explicitly if needed
set -euxo pipefail

# Clone libaom-v1
if [ ! -d "./av1" ]; then
      git clone --branch master --depth 1 https://aomedia.googlesource.com/aom av1
fi

# Build libaom-v1
mkdir -p ./av1/build && cd ./av1/build
cmake -G "Unix Makefiles" -DCMAKE_INSTALL_PREFIX="/usr/local" -DENABLE_SHARED=off -DENABLE_NASM=on ..
make -j $(nproc)

# Install libaom-v1
make install