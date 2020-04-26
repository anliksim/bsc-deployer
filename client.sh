#!/bin/bash

to_json() {
  cat <<EOF
{
  "rev": "$1",
  "dir": "$2"
}
EOF
}

case "$1" in
"stop")
  http_method=DELETE
  ;;
*)
  http_method=POST
  ;;
esac

headRev=$(git log --pretty=format:'%h %s' --abbrev-commit -1)
workdir=$(dirname "$(pwd)")
data=$(to_json "$headRev" "$workdir/bsc-env")

out=$(curl -k -sS -X $http_method "http://localhost:3557/v1/deployments" --data "$data")
echo "Deployment started at $out"
