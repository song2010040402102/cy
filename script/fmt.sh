#!/bin/bash

cd ../src
go fmt
for dir in $(ls ./)
do
	if [ -d $dir ]; then
		cd $dir && go fmt
		cd ..
	fi
done  