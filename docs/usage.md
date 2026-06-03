<!-- This file is auto-generated. Do not edit manually. -->

# zxcv

zxcv - A minimalist, fast .tool-versions manager

## DESCRIPTION

zxcv is a dead simple tool-version manager — a cleaner, friendlier alternative to asdf that plugs into its ecosystem with a highly minimalistic command surface and enjoyably better performance.

**Usage**:

```
zxcv [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

## COMMANDS

### current

Show resolved tools and their versions for the current directory

**--version**: output version(s) only

### search

Search the upstream tool registry by tool name, or a tool's available versions

### install

Install resolved tools or add new ones

**--global**: install globally; affects the global ~/.tool-versions file

**--no-save**: do not modify any .tool-versions file

### exec

Run a specific binary

### which

Output the resolved path of a tool's binary

### list

List installed tools and their versions

### uninstall

Remove installed tools
