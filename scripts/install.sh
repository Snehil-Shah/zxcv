#!/usr/bin/env sh
# zxcv installer.
#
# Usage:
#   curl -fsSL install.sh | sh
#   curl -fsSL install.sh | sh -s v0.1.0
#
# Env overrides:
#   ZXCV_VERSION       pin a release tag (e.g. v0.1.0). Default: latest.
#   ZXCV_INSTALL_DIR   binary install dir. Default: ~/.local/bin.
#   ZXCV_NO_MODIFY_RC  if 1, skip shell rc modification (PATH + completions).

set -eu

REPO="Snehil-Shah/zxcv"
VERSION="${1:-${ZXCV_VERSION:-}}"
INSTALL_DIR="${ZXCV_INSTALL_DIR:-$HOME/.local/bin}"
SHIMS_DIR="$HOME/.zxcv/bin"

err()  { printf '\033[31merror:\033[0m %s\n' "$1" >&2; exit 1; }
info() { printf '\033[2m%s\033[0m\n' "$1"; }
ok()   { printf '\033[32m✓\033[0m %s\n' "$1"; }

# Normalize "0.1.0" -> "v0.1.0"; leave empty and "v*" as-is.
case "$VERSION" in
  ""|v*) ;;
  *) VERSION="v${VERSION}" ;;
esac

# Detect OS.
case "$(uname -s)" in
  Darwin) os="darwin" ;;
  Linux)  os="linux" ;;
  *) err "unsupported OS: $(uname -s)" ;;
esac

# Detect arch.
case "$(uname -m)" in
  x86_64|amd64)  arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) err "unsupported arch: $(uname -m)" ;;
esac

# Resolve latest if no version was specified.
if [ -z "$VERSION" ]; then
  info "resolving latest release..."
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name":' | head -n1 | sed -E 's/.*"([^"]+)".*/\1/')"
  [ -n "$VERSION" ] || err "could not resolve latest release"
fi

asset="zxcv_${os}_${arch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"

info "downloading ${asset} (${VERSION})..."
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

if ! curl -fsSL "$url" -o "${tmpdir}/${asset}"; then
  err "no release ${VERSION} for ${os}/${arch} — see https://github.com/${REPO}/releases for available versions"
fi
tar -xzf "${tmpdir}/${asset}" -C "$tmpdir" || err "extract failed"

mkdir -p "$INSTALL_DIR"
mv "${tmpdir}/zxcv" "${INSTALL_DIR}/zxcv"
chmod +x "${INSTALL_DIR}/zxcv"
ok "installed zxcv ${VERSION} to ${INSTALL_DIR}/zxcv"

# True iff $1 is a path segment in $PATH.
on_path() {
  case ":${PATH}:" in *":$1:"*) return 0 ;; esac
  return 1
}

# Convert "$HOME/foo" -> literal "$HOME/foo" string for rc files (portable across home moves).
to_rc_path() {
  case "$1" in
    "$HOME"/*) printf '$HOME/%s' "${1#"$HOME"/}" ;;
    *) printf '%s' "$1" ;;
  esac
}

# Decide which dirs need PATH entries.
add_install=0; on_path "$INSTALL_DIR" || add_install=1
add_shims=0;   on_path "$SHIMS_DIR"   || add_shims=1
install_dir_rc="$(to_rc_path "$INSTALL_DIR")"
shims_dir_rc="$(to_rc_path "$SHIMS_DIR")"

# Append a managed PATH + completion block to rc file (idempotent).
modify_rc() {
  rc="$1"
  shell="$2"  # bash, zsh, or fish
  if [ -f "$rc" ] && grep -qF "# >>> zxcv >>>" "$rc"; then
    return 0
  fi
  mkdir -p "$(dirname "$rc")"
  {
    printf '\n# >>> zxcv >>>\n'
    if [ "$shell" = "fish" ]; then
      [ "$add_install" = "1" ] && printf 'set -gx PATH "%s" $PATH\n' "$install_dir_rc"
      [ "$add_shims"   = "1" ] && printf 'set -gx PATH "%s" $PATH\n' "$shims_dir_rc"
    else
      [ "$add_install" = "1" ] && printf 'export PATH="%s:$PATH"\n' "$install_dir_rc"
      [ "$add_shims"   = "1" ] && printf 'export PATH="%s:$PATH"\n' "$shims_dir_rc"
      printf 'source <(zxcv completion %s)\n' "$shell"
    fi
    printf '# <<< zxcv <<<\n'
  } >> "$rc"
  ok "wired ${rc}"
  rc_modified=1
}

# Install fish completions as a discoverable file (fish auto-loads from this dir).
install_fish_completions() {
  fish_comp_dir="$HOME/.config/fish/completions"
  fish_comp_file="${fish_comp_dir}/zxcv.fish"
  mkdir -p "$fish_comp_dir"
  "${INSTALL_DIR}/zxcv" completion fish > "$fish_comp_file"
  ok "wrote fish completions to ${fish_comp_file}"
}

rc_modified=0

if [ "${ZXCV_NO_MODIFY_RC:-}" = "1" ]; then
  info "skipped shell rc modification (ZXCV_NO_MODIFY_RC=1)"
else
  case "$(basename "${SHELL:-sh}")" in
    bash)
      modify_rc "$HOME/.bashrc" bash
      [ -f "$HOME/.bash_profile" ] && modify_rc "$HOME/.bash_profile" bash
      ;;
    zsh)
      modify_rc "$HOME/.zshrc" zsh
      ;;
    fish)
      modify_rc "$HOME/.config/fish/config.fish" fish
      install_fish_completions
      ;;
    *)
      info "unknown shell '${SHELL:-sh}' — add to your shell rc manually:"
      [ "$add_install" = "1" ] && info "  export PATH=\"${install_dir_rc}:\$PATH\""
      [ "$add_shims"   = "1" ] && info "  export PATH=\"${shims_dir_rc}:\$PATH\""
      info "  # plus your shell's equivalent of: source <(zxcv completion <shell>)"
      ;;
  esac
fi

if [ "$rc_modified" = "1" ]; then
  printf '\nrestart your shell or `source` your rc file to pick up changes.\n'
fi
