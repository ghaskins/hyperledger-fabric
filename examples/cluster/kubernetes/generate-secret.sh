#!/bin/bash
#
# Copyright Greg Haskins All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

NAME=$1
FILE=$2

SCRIPT_DIR=$(dirname $0)

cat <<EOF
---
apiVersion: v1
kind: Secret
metadata:
  name: $NAME
type: Opaque
data:
  config.tgz: $(cat $FILE | base64)
EOF
