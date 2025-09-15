#####################################
FROM nvidia/cuda:11.4.3-devel-ubuntu20.04 AS ubi-builder

# Prevent interactive prompts from apt
ENV DEBIAN_FRONTEND=noninteractive
 
# Set environment variables for Go version and architecture
# This makes it easy to update Go later
ENV GO_VERSION="1.22.12"
ENV GO_OS="linux"
ENV GO_ARCH="amd64"
# SHA256 checksum for go1.22.0.linux-amd64.tar.gz - ALWAYS verify this from golang.org/dl
# You can find the correct checksum on the official Go downloads page:
# https://golang.org/dl/
ENV GO_SUM="4fa4f869b0f7fc6bb1eb2660e74657fbf04cdd290b5aef905585c86051b34d43"
 
# Install necessary packages for downloading and extracting Go, and common build tools
# --no-install-recommends helps keep the image size down
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl wget git jq pkg-config hwloc libhwloc-dev coreutils vim && \
    rm -rf /var/lib/apt/lists/*
 
# Download, verify, and install Go
# -fsSL: Fail silently, show errors, follow redirects
# The checksum verification is crucial for security and integrity.
RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.${GO_OS}-${GO_ARCH}.tar.gz" -o go.tar.gz && \
    echo "${GO_SUM} go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz
 
# Set Go environment variables
ENV GOROOT=/usr/local/go
# Add GOROOT/bin and GOPATH/bin to the system PATH
ENV PATH=$PATH:$GOROOT/bin
ENV GOPATH=/go
 
# Create GOPATH directory
RUN mkdir -p ${GOPATH}

ENV XDG_CACHE_HOME="/tmp"

### taken from https://github.com/rust-lang/docker-rust/blob/master/1.63.0/buster/Dockerfile
ENV RUSTUP_HOME=/usr/local/rustup \
    CARGO_HOME=/usr/local/cargo \
    PATH=/usr/local/cargo/bin:$PATH \
    RUST_VERSION=1.63.0

RUN set -eux; \
    dpkgArch="$(dpkg --print-architecture)"; \
    case "${dpkgArch##*-}" in \
        amd64) rustArch='x86_64-unknown-linux-gnu'; rustupSha256='5cc9ffd1026e82e7fb2eec2121ad71f4b0f044e88bca39207b3f6b769aaa799c' ;; \
        arm64) rustArch='aarch64-unknown-linux-gnu'; rustupSha256='e189948e396d47254103a49c987e7fb0e5dd8e34b200aa4481ecc4b8e41fb929' ;; \
        *) echo >&2 "unsupported architecture: ${dpkgArch}"; exit 1 ;; \
    esac; \
    url="https://static.rust-lang.org/rustup/archive/1.25.1/${rustArch}/rustup-init"; \
    wget "$url"; \
    echo "${rustupSha256} *rustup-init" | sha256sum -c -; \
    chmod +x rustup-init; \
    ./rustup-init -y --no-modify-path --profile minimal --default-toolchain $RUST_VERSION --default-host ${rustArch}; \
    rm rustup-init; \
    chmod -R a+w $RUSTUP_HOME $CARGO_HOME; \
    rustup --version; \
    cargo --version; \
    rustc --version;

COPY ./ /opt/ubi-benchmark
WORKDIR /opt/ubi-benchmark

RUN mkdir -p /usr/lib/x86_64-linux-gnu/
# NOTE: copy your host /usr/lib/x86_64-linux-gnu/libcuda.so* to the libcuda_files folder first
COPY libcuda_files/* /usr/lib/x86_64-linux-gnu/

ENV LD_LIBRARY_PATH=/usr/lib/x86_64-linux-gnu:/usr/local/cuda/lib64:${LD_LIBRARY_PATH}

### make configurable filecoin-ffi build
ARG FFI_BUILD_FROM_SOURCE=1
ENV FFI_BUILD_FROM_SOURCE=${FFI_BUILD_FROM_SOURCE}
ENV RUSTFLAGS="-C target-cpu=native -g"
ENV FFI_USE_CUDA=1

RUN make clean build

#####################################
FROM ubuntu:20.04 AS ubi-benchmark

ENV DEBIAN_FRONTEND=noninteractive
RUN ln -fs /usr/share/zoneinfo/America/Toronto /etc/localtime \
    && apt-get update \
    && apt-get install -y tzdata \
    && dpkg-reconfigure -f noninteractive tzdata

COPY --from=ubi-builder /opt/ubi-benchmark/ubi-bench /usr/local/bin/
ENV TRUST_PARAMS=1
ENV RUST_LOG=Info
ENV UBI_TASK_IN_PARAM_PATH=/var/tmp/fil-c2-param
ENV FILECOIN_PARAMETER_CACHE=/var/tmp/filecoin-proof-parameters

RUN apt-get update && apt-get install -y hwloc libhwloc-dev coreutils vim
RUN mkdir /var/tmp/filecoin-proof-parameters
VOLUME /var/tmp/filecoin-proof-parameters

CMD ["/bin/bash", "-c", "sleep infinity"]
