#!/bin/bash
./regression.py
clc -s -Lpy
cat Version.dat
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st
