#!/usr/bin/env sh
# zxcv uninstaller.
#
# Usage:
#   curl -fsSL uninstall.sh | sh
#
# Env overrides:
#   ZXCV_INSTALL_DIR   binary location. Default: ~/.local/bin.
#   ZXCV_NO_MODIFY_RC  if 1, skip cleanup of shell rc files (PATH + completions).

set -eu

INSTALL_DIR="${ZXCV_INSTALL_DIR:-$HOME/.local/bin}"
DATA_DIR="$HOME/.zxcv"

err()  { printf '\033[31merror:\033[0m %s\n' "$1" >&2; exit 1; }
info() { printf '\033[2m%s\033[0m\n' "$1"; }
ok()   { printf '\033[32m✓\033[0m %s\n' "$1"; }

binary="${INSTALL_DIR}/zxcv"

# Run zxcv's own uninstaller first.
if [ -x "$binary" ]; then
  info "running zxcv uninstall..."
  if ! "$binary" uninstall; then
    info "zxcv uninstall failed — continuing with raw cleanup"
  fi
fi

# Remove the binary.
if [ -e "$binary" ] || [ -L "$binary" ]; then
  rm -f "$binary"
  ok "removed ${binary}"
fi

# Strip the managed PATH block from a shell rc file (idempotent).
remove_rc_block() {
  rc="$1"
  [ -f "$rc" ] || return 0
  grep -qF "# >>> zxcv >>>" "$rc" || return 0
  tmp="$(mktemp)"
  sed '/# >>> zxcv >>>/,/# <<< zxcv <<</d' "$rc" > "$tmp" && mv "$tmp" "$rc"
  ok "removed PATH block from ${rc}"
}

if [ "${ZXCV_NO_MODIFY_RC:-}" = "1" ]; then
  info "skipped rc cleanup (ZXCV_NO_MODIFY_RC=1)"
else
  for rc in \
    "$HOME/.bashrc" \
    "$HOME/.bash_profile" \
    "$HOME/.zshrc" \
    "$HOME/.config/fish/config.fish"
  do
    remove_rc_block "$rc"
  done
  # Fish completions live in their own file outside the managed block.
  fish_comp="$HOME/.config/fish/completions/zxcv.fish"
  if [ -f "$fish_comp" ]; then
    rm -f "$fish_comp"
    ok "removed ${fish_comp}"
  fi
fi

# Wipe whatever's left of the data dir.
if [ -d "$DATA_DIR" ]; then
  rm -rf "$DATA_DIR"
  ok "removed ${DATA_DIR}"
fi
