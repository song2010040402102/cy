#!/bin/bash

if [ ! -d "../bin" ]; then
  mkdir ../bin
fi

export GOPATH="$GOPATH:$(cd ..; pwd)"

cd ../src
go build -o ../bin/cy