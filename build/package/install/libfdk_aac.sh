#!/bin/bash

# Fail explicitly if needed
set -euxo pipefail

# Clone libfdk_aac
if [ ! -d "./fdk-aac" ]; then
      git clone --depth 1 git://git.code.sf.net/p/opencore-amr/fdk-aac
fi

# Build libfdk_aac
cd fdk-aac
autoreconf -fiv
./configure --prefix="/usr/local" --disable-shared
make -j $(nproc)

# Install libfdk_aac
make install