#!/bin/sh
set -e

REPO="chainreactors/libcstx"
VERSION="${CSTX_VERSION:-latest}"

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case "$OS" in
        linux)  OS="linux" ;;
        darwin) OS="darwin" ;;
        mingw*|msys*|cygwin*) OS="windows" ;;
        *) echo "unsupported OS: $OS" >&2; exit 1 ;;
    esac
    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "unsupported arch: $ARCH" >&2; exit 1 ;;
    esac
    echo "${OS}_${ARCH}"
}

PLATFORM=$(detect_platform)
LIB_DIR="$(cd "$(dirname "$0")/.." && pwd)/lib/${PLATFORM}"
MARKER="${LIB_DIR}/.version"

if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "failed to fetch latest release version" >&2
        exit 1
    fi
fi

if [ -f "$MARKER" ] && [ "$(cat "$MARKER")" = "$VERSION" ]; then
    echo "libcstx_ffi ${VERSION} already installed for ${PLATFORM}"
    exit 0
fi

echo "downloading libcstx_ffi ${VERSION} for ${PLATFORM}..."
URL="https://github.com/${REPO}/releases/download/${VERSION}/libcstx_ffi-${PLATFORM}.tar.gz"

mkdir -p "$LIB_DIR"
curl -sL "$URL" | tar -xz -C "$LIB_DIR"
echo "$VERSION" > "$MARKER"
echo "installed libcstx_ffi ${VERSION} → ${LIB_DIR}"
