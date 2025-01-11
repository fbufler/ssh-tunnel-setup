#!/bin/bash

# Build the go binary based on the architecture argument passed

arch=$1
os=$2

script_dir=$(dirname $0)

if [ -z "$arch" ]; then
  echo "Please provide the architecture as an argument"
  exit 1
fi

if [ -z "$os" ]; then
  echo "Please provide the os as an argument"
  exit 1
fi
if [ "$os" == "darwin" ]; then
  if [ "$arch" == "amd64" ]; then
    GOOS=darwin GOARCH=amd64 go build -o ${script_dir}/../bin/ssh-tunnel-setup-darwin-amd64 ${script_dir}/../main.go
  elif [ "$arch" == "arm64" ]; then
    GOOS=darwin GOARCH=arm64 go build -o ${script_dir}/../bin/ssh-tunnel-setup-darwin-arm64 ${script_dir}/../main.go
  else
    echo "Invalid architecture"
    exit 1
  fi
elif [ "$os" == "linux" ]; then
  if [ "$arch" == "amd64" ]; then
    GOOS=linux GOARCH=amd64 go build -o ${script_dir}/../bin/ssh-tunnel-setup-linux-amd64 ${script_dir}/../main.go
  elif [ "$arch" == "arm64" ]; then
    GOOS=linux GOARCH=arm64 go build -o ${script_dir}/../bin/ssh-tunnel-setup-linux-arm64 ${script_dir}/../main.go
  else
    echo "Invalid architecture"
    exit 1
  fi
else
  echo "Invalid os"
  exit 1
fi