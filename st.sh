#!/bin/bash
./regression.py -q
clc -s -Lpy
cat Version.dat
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st
