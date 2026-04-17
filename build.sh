#!/bin/bash
mkdir build
go build shex.go
cd compiler
go build compile.go
mv ../shex ../build
mv compile ../build