# syntax=docker/dockerfile:1

ARG GO_VERSION=1.26.4
ARG BUILD_DIR=/build

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS builder-base
WORKDIR /src

FROM builder-base AS builder-deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    go mod download

FROM builder-deps AS builder-lint
RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    make lint

FROM builder-lint AS builder-test
ARG BUILD_DIR
RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    mkdir -p ${BUILD_DIR} && \
    make test BUILD_DIR=${BUILD_DIR}

FROM builder-test AS builder
ARG TARGETOS
ARG TARGETARCH
ARG BUILD_DIR
RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    mkdir -p ${BUILD_DIR} && \
    make build GOOS=${TARGETOS} GOARCH=${TARGETARCH} BUILD_DIR=${BUILD_DIR}

FROM scratch AS runtime
ARG BUILD_DIR
WORKDIR /
COPY --from=builder ${BUILD_DIR}/bin/cdrdao2audio /cdrdao2audio
ENTRYPOINT ["/cdrdao2audio"]
