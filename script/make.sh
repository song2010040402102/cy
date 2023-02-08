#!/bin/bash

cd $(dirname $0)
cd ..

if [ ! -d "bin" ]; then
  mkdir bin
fi

go build -o ./bin/cy