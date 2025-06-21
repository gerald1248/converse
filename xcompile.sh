#!/bin/bash

FILENAME=converse

rm -f release.txt

mkdir -p dist/
for OS in windows linux darwin; do
  mkdir -p ${OS}
  GOOS=${OS} GOARCH=amd64 go build -o ${OS}/${FILENAME}
  if [[ "$OS" == "windows" ]]; then
    mv windows/${FILENAME} windows/${FILENAME}.exe
  fi
  zip -jr dist/${FILENAME}-${OS}-amd64.zip ${OS}/${FILENAME}*
  rm -rf ${OS}/
  shasum -a 256 dist/${FILENAME}-${OS}-amd64.zip >>dist/release.txt
done
