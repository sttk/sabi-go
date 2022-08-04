#!/usr/bin/env bash

clean() {
  go clean --cache
}

format() {
  go fmt ./...
}

compile() {
  go vet ./...
  go build ./...
}

test() {
  go test -v ./...
}

cover() {
  mkdir -p coverage
  go test -coverprofile=coverage/cover.out ./...
  go tool cover -html=coverage/cover.out -o coverage/cover.html
}

bench() {
  local dir=$2
  if [[ "$dir" == "" ]]; then
    dir="."
  fi
  pushd $dir
  go test -bench . -benchmem -run=^$
  popd
}

if [[ $# == 0 ]]; then
  clean
  format
  compile
  test
  cover
  exit 0
fi

for a in "$@"; do
  case "$a" in
  clean)
    clean
    ;;
  format)
    format
    ;;
  compile)
    compile
    ;;
  test)
    test
    ;;
  cover)
    cover
    ;;
  bench)
    bench
    ;;
  '')
    compile
    ;;
  *)
    echo "Bad task: $a"
    exit 1
    ;;
  esac
done
