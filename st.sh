#!/bin/bash
clc -s -e murmur_test.go
cat Version.dat
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st
