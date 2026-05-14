# .tool-definitions

Drop a `.tool-definitions` file alongside your `.tool-versions` to override or extend the upstream tools registry, for tools that aren't in the registry.

This is what a .tool-definitions file looks like:

```
internal-cli https://github.com/our-org/asdf-internal-cli.git
my-fork-of-nodejs git@github.com:me/asdf-nodejs-fork.git
my-plugin ./vendor/my-plugin
```

Each line is `<tool> <source>`. The source can be a `git://` or `https://` URL, an SSH-style git URL, or a filesystem path (absolute, or relative to the `.tool-definitions` file).

You can also include comments:

```
# Internal tools
internal-cli https://github.com/our-org/asdf-internal-cli.git
# This is another comment
my-plugin ./vendor/my-plugin
```

!!! tip
    Filesystem paths are symlinked instead of cloned. Handy when authoring plugins — edit your local checkout in place and re-install to see changes.
