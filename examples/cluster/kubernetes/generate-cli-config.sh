#!/bin/bash
#
# Copyright Greg Haskins All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

ORDERER_CRT=$1

SCRIPT_DIR=$(dirname $0)

includefile() {
    file=$1
    prefix="    "

    echo "|"

    while IFS= read -r line; do
        printf '%s%s\n' "$prefix" "$line"
    done < "$file"
}

cat <<EOF
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: configurator-payload
data:
  orderer-ca.crt: $(includefile $ORDERER_CRT)
  configure.sh: $(includefile $SCRIPT_DIR/configure.sh.in)
EOF
