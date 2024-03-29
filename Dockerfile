#####################################
FROM sxk1633/lotus-compile-cuda:11.4 AS ubi-builder
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


### make configurable filecoin-ffi build
ARG FFI_BUILD_FROM_SOURCE=1
ENV FFI_BUILD_FROM_SOURCE=${FFI_BUILD_FROM_SOURCE}
ENV RUSTFLAGS="-C target-cpu=native -g"
ENV FFI_USE_CUDA=1

RUN make clean build

#####################################
FROM ubuntu:20.04 AS ubi-benchmark

COPY --from=ubi-builder /opt/ubi-benchmark/ubi-bench /usr/local/bin/
ENV TRUST_PARAMS=1
ENV RUST_LOG=Info
ENV UBI_TASK_IN_PARAM_PATH /var/tmp/fil-c2-param
ENV FILECOIN_PARAMETER_CACHE /var/tmp/filecoin-proof-parameters

RUN apt-get update && apt-get install -y hwloc libhwloc-dev coreutils vim
RUN mkdir /var/tmp/filecoin-proof-parameters
RUN chown fc: /var/tmp/filecoin-proof-parameters

VOLUME /var/tmp/filecoin-proof-parameters

CMD ["/bin/bash", "-c", "sleep infinity"]
