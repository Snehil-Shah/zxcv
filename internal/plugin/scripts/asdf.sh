#!/usr/bin/env sh
# Minimal asdf shim.
# HACK: This shim translates some common asdf invocations (generally by plugin scripts) to its zxcv alternatives

set -eu

case "$1" in
  current)
    shift
    tool=""
    while [ "$#" -gt 0 ]; do
      case "$1" in
        --no-header) ;;
        --*) ;;
        *) tool="$1" ;;
      esac
      shift
    done
    if [ -z "$tool" ]; then
      echo "asdf-shim: 'current' requires a tool" >&2
      exit 1
    fi
    version=$(zxcv current "$tool" --version)
    echo "$tool $version"
    ;;
  version)
    echo "v0.19.0"  # HACK: This is just so the command doesn't break the parent script
    ;;
  reshim)
    exit 0
    ;;
  *)
    echo "asdf-shim: '$1' not supported under zxcv" >&2
    exit 1
    ;;
esac
