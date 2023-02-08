#!/bin/bash

cd $(dirname $0)

fmt_dir(){
    for file in `ls $1`
    do
        if [ -d $1"/"$file ]; then
            fmt_dir $1"/"$file
        else
            if [[ $file == *.go ]]; then
                go fmt $1"/"$file
            fi
        fi
    done
}
fmt_dir ..