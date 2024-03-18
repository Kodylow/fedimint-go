#!/bin/bash

current_dir=$(pwd)
cd $current_dir/pkg/fedimint
go test
cd "$current_dir"
