# .tool-versions

Whenever a `.tool-versions` file is present in a directory, the tool versions it declares will be used in that directory and any subdirectories.

This is what a .tool-versions file looks like:

```
ruby 2.5.3
nodejs 10.15.0
```

You can also include comments:

```
ruby 2.5.3 # This is a comment
# This is another comment
nodejs 10.15.0
```

!!! tip
    Use the [`search`](../reference/search.md) command to find the exact tool names and versions you wish to use.

You can specify `system` in place of the version if you wish to use the system installed binary instead of `zxcv` managed ones.

!!! note
    The entire [asdf-spec](https://asdf-vm.com/manage/configuration.html#tool-versions) is not supported. Only a simpler version of it defined above.