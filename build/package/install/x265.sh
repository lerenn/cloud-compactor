#!/bin/bash

# Fail explicitly if needed
set -euxo pipefail

# Clone x265
if [ ! -d "./x265" ]; then
      git clone https://bitbucket.org/multicoreware/x265_git x265
fi

# Build x265
mkdir -p ./x265/build/linux && cd ./x265/build/linux
cmake -G "Unix Makefiles" -DCMAKE_INSTALL_PREFIX="/usr/local" -DENABLE_SHARED:bool=off ../../source
make -j $(nproc)

# Install x265
make install