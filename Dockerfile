FROM debian:stable-slim AS builder

ARG TARGETOS
ARG TARGETARCH
ARG BIN_NAME
ARG TAGS

RUN apt update && apt install -y --no-install-recommends llvm-19-dev wget gcc g++ ca-certificates
RUN wget -c -O go.tar.gz https://go.dev/dl/go1.26.5.$TARGETOS-$TARGETARCH.tar.gz
RUN tar -C /usr/local -xzf go.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

ENV CGO_ENABLED=1 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build $TAGS -v -o /out/$BIN_NAME .

FROM scratch
ARG BIN_NAME

COPY --from=builder /out/$BIN_NAME /$BIN_NAME
