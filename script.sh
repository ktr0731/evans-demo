#!/usr/bin/env bash

set -e

request=$(cat <<EOS
{
  "users": [
    {
      "first_name": "tsukasa",
      "last_name": "ayatsuji",
      "gender": 1
    },
    {
      "first_name": "kaoru",
      "last_name": "tanamachi",
      "gender": 1
    },
    {
      "first_name": "haruka",
      "last_name": "morishima",
      "gender": 1
    }
  ]
}
EOS
)

echo  "$request" | $GOPATH/bin/evans --call RegisterUsers > /dev/null
echo '{}' | $GOPATH/bin/evans --call ListUsers















