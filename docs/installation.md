# Installation

## Quick install

```sh
curl -fsSL https://snehil-shah.github.io/zxcv/install.sh | sh
```

This downloads the latest release for your OS/architecture and also modifies you shell rc (see customization below if you wish to avoid this).

After install, restart your shell or `source` your rc file.

## Pinning a version

```sh
# Positional argument:
curl -fsSL https://snehil-shah.github.io/zxcv/install.sh | sh -s v0.1.0

# Or via environment variable:
ZXCV_VERSION=v0.1.0 curl -fsSL https://snehil-shah.github.io/zxcv/install.sh | sh
```

The `v` prefix is auto-added if you forget it (`0.1.0` works the same).

## Customizing installation

Set environment variables before running the script:

| Variable | Default | Effect |
|---|---|---|
| `ZXCV_VERSION` | (latest) | Pin a specific release tag. |
| `ZXCV_INSTALL_DIR` | `~/.local/bin` | Where to place the `zxcv` binary. |
| `ZXCV_NO_MODIFY_RC` | (unset) | If `1`, skip shell rc file modifications (PATH + completions). |

## Uninstalling

```sh
curl -fsSL https://snehil-shah.github.io/zxcv/uninstall.sh | sh
```

## Supported platforms

- macOS (arm64, amd64)
- Linux (arm64, amd64)

Windows is not supported. WSL works since it presents as Linux.
