#!/bin/bash

OUTPUT=$1
INCLUDES=$2

echo "Generating $OUTPUT with [$INCLUDES]"

cat <<EOF > $OUTPUT
package syscc

import (
`for i in $INCLUDES; do echo "\"$i\""; done`
)
EOF

gofmt -w $OUTPUT
