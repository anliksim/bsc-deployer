#!/bin/bash

to_json() {
  cat <<EOF
{
  "rev": "$1",
  "dir": "$2"
}
EOF
}

headRev=$(git log --pretty=format:'%h %s' --abbrev-commit -1)
workdir=$(dirname `pwd`)

data=$(to_json "$headRev" "$workdir/bsc-env/apps")

out=$(curl -k -sS -X POST "http://localhost:3557/v1/deployments" --data "$data")
#curl -k -X DELETE "http://localhost:3557/v1/deployments" --data "$data"
echo "Deployment started at $out"

