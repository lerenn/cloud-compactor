#!/bin/bash

# Fail explicitly if needed
set -euxo pipefail

# Clone libmp3lame
if [ ! -d "./lame-3.100" ]; then
      curl -L -O http://downloads.sourceforge.net/project/lame/lame/3.100/lame-3.100.tar.gz
      tar xzvf lame-3.100.tar.gz
fi

# Build libmp3lame
cd lame-3.100
autoreconf --install
./configure --prefix="/usr/local" --disable-shared --enable-nasm
make -j $(nproc)

# Install libmp3lame
make install