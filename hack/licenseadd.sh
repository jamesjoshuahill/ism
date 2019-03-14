#!/bin/bash

for i in $(find . -name '*.go' -not -path "./vendor/*");
do
  if ! grep -q Copyright "$i"
  then
    cat ./hack/boilerplate.go.txt "$i" >"$i".new && mv "$i".new "$i"
  fi
done

for i in $(find . -name '*.yaml' -o -name '*.yml' -not -path "./vendor/*");
do
  if ! grep -q Copyright "$i"
  then
    cat ./hack/boilerplate.yaml.txt "$i" >"$i".new && mv "$i".new "$i"
  fi
done
