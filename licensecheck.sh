#!/bin/bash

for i in $(find . -type f -not -path "./vendor/*" -not -path "./.git/*"); do
  if ! grep -q Copyright "$i"; then
    echo "$i"
  fi
done