FROM fedora:43

# these don't depend on git sha, and are pretty much the same on every build
RUN dnf install -y \
    gcc-c++ \
    libcxx-devel \
    openssl-devel \
    zlib-ng-compat-devel \
    gperf \
    cmake \
    git

# WORKDIR /

# https://tdlib.github.io/td/build.html?language=Go
RUN git clone --branch master \
    --single-branch \
    --depth 1 \
    https://github.com/tdlib/td.git \
    td/

RUN cd td \
    && rm -rf build \
    && mkdir build \
    && cd build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=../tdlib .. \
    && cmake --build . --target install \
    && cd ../.. \
    && ls -l td/tdlib
