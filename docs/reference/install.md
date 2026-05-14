<!-- This file is auto-generated. Do not edit manually. -->

# install

install - Install resolved tools or add new ones

## DESCRIPTION

Downloads & installs all resolved tools for the current directory, if no arguments are provided. If a specific tool and its version are provided, it installs only that and adds them to the local .tool-versions manifest. This behavior can be controlled to save it to the global manifest or not save at all using the --global and --no-save flags, respectively. You can specify the version as 'latest' (or with a prefix 'latest:<prefix>') to simply install and pin the latest version of the tool.

**Usage**:

```
install [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

## GLOBAL OPTIONS

**--global**: install globally; affects the global ~/.tool-versions file

**--no-save**: do not modify any .tool-versions file

