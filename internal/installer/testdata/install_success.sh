#!/bin/sh
mkdir -p "$ASDF_INSTALL_PATH/bin"
echo "version $ASDF_INSTALL_VERSION" > "$ASDF_INSTALL_PATH/bin/mytool"
chmod +x "$ASDF_INSTALL_PATH/bin/mytool"
