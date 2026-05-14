# Getting Started

This page walks through your first install with zxcv. Assumes you've already [installed zxcv](installation.md).

## Installing tools from `.tool-versions`

zxcv reads the standard asdf [`.tool-versions` format](./spec/tool-versions.md). Drop a file like this in your project:

```
nodejs 20.11.0
python 3.11.5
ruby 3.2.0
```

Then:

```sh
zxcv install
```

## Installing a single tool

```sh
zxcv install nodejs 20.11.0
zxcv install python latest
zxcv install golang latest:1.26
```

By default this also pins the version into your local `.tool-versions`. Use `--no-save` to install without pinning, or `--global` to write to `$HOME/.tool-versions` instead.

## Discovering available tools and versions you can install

```sh
zxcv search node          # tools matching "node"
zxcv search nodejs latest # latest stable version of nodejs
zxcv search nodejs 20     # all nodejs versions starting with "20"
```

## Running installed binaries

Once installed, the binaries are on your `PATH`. Just use them directly:

```sh
node --version
python --version
```

!!! tip
    If you are confused on what all tools are resolved for your directory, run the [`zxcv current`](reference/current.md) command.

To inspect which version our binary resolves to:

```sh
zxcv which node
```

To run a binary at a specific version without changing the active version:

```sh
zxcv exec nodejs 18.19.0 -- node script.js
```
