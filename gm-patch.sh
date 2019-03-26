#!/bin/bash

sed -i 's/mediocregopher\/radix\/v3/groupme\/radix/g' *.go
sed -i 's/mediocregopher\/radix\/v3/groupme\/radix/g' ./resp/resp2/*.go
rm ./go.mod
rm ./go.sum
rm -rf ./bench

# Because those depends on stretchr libraries which fails
rm ./scanner_test.go
rm ./stream_test.go

echo "# Radix

forked and patched from github.com/mediocregopher/radix to overcome issues with golang modules

" > ./README.md