#!/bin/sh
set -eu

REPOSITORY="duck-ahiru-Z/DevDuck"
INSTALL_DIR="${DEVDUCK_INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${DEVDUCK_VERSION:-}"

OS=$(uname -s)
ARCH=$(uname -m)
case "$OS:$ARCH" in
  Linux:x86_64|Linux:amd64) ASSET="devduck-linux-amd64" ;;
  Linux:aarch64|Linux:arm64) ASSET="devduck-linux-arm64" ;;
  Darwin:x86_64|Darwin:amd64) ASSET="devduck-darwin-amd64" ;;
  Darwin:arm64) ASSET="devduck-darwin-arm64" ;;
  *) echo "Unsupported platform: $OS/$ARCH" >&2; exit 1 ;;
esac

if [ -n "$VERSION" ]; then
  BASE_URL="https://github.com/$REPOSITORY/releases/download/$VERSION"
else
  BASE_URL="https://github.com/$REPOSITORY/releases/latest/download"
fi

tmpdir=$(mktemp -d 2>/dev/null || mktemp -d -t devduck-install)
cleanup() { rm -rf "$tmpdir"; }
trap cleanup EXIT INT TERM

binary="$tmpdir/$ASSET"
checksums="$tmpdir/SHA256SUMS"
curl -fsSL "$BASE_URL/$ASSET" -o "$binary"
curl -fsSL "$BASE_URL/SHA256SUMS" -o "$checksums"

expected=$(awk -v asset="$ASSET" '$2 == asset {print $1; exit}' "$checksums")
[ -n "$expected" ] || { echo "No checksum found for $ASSET" >&2; exit 1; }
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$binary" | awk '{print $1}')
else
  actual=$(shasum -a 256 "$binary" | awk '{print $1}')
fi
[ "$actual" = "$expected" ] || { echo "SHA256 mismatch for $ASSET" >&2; exit 1; }

mkdir -p "$INSTALL_DIR"
chmod +x "$binary"
mv "$binary" "$INSTALL_DIR/duck"
echo "Installed DevDuck at $INSTALL_DIR/duck"
case ":${PATH:-}:" in
  *:"$INSTALL_DIR":*) ;;
  *) echo "$INSTALL_DIR is not on PATH. Add it to your shell configuration, then open a new shell." ;;
esac
