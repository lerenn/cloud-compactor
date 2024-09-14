#!/bin/bash

# Fail explicitly if needed
set -euxo pipefail

# Clone x264
if [ ! -d "./x264" ]; then
    git clone --depth 1 https://code.videolan.org/videolan/x264.git;
fi 

# Build x264
cd x264
./configure --prefix="/usr/local" --enable-static 
make -j $(nproc)

# Install x264
make install 