#!/bin/bash

for i in $(find . -type f -not -path "./vendor/*" -not -path "./.git/*" -not -path "./removelicense.sh"); do
  if grep -q "Copyright 2019 Independent Services Marketplace Team" "$i"; then
    sed '1,16d' "$i" > "$i".new && mv "$i".new "$i"
  fi
done
