#!/bin/sh
# This script will create a completely static build of the program,
# or any Go program you use it with, so far as I am aware.
# So long as you build with 'go build .',
# this script should work correctly.
# Modification is outside the scope of the comments.
# To use this script, you need musl libc installed.
# You will also need musl-gcc in your path, so the script can properly execute.
# This will link your program against the musl lib and create a completely static program.
# Not all Go programs will do this by default, but if you plan to use this in a Docker container such as scratch,
# you will need a completely static build.
# In addition to making a static build, this script will strip all debug information from the binaries,
# and trim the build path of your current working directory from all information returned in case of a panic.

CC=musl-gcc go build -buildmode=pie "-asmflags=all='-trimpath=`pwd`'" -ldflags "-linkmode external -w -s -extldflags '-static' -X main.Version=$(git describe --tags `git rev-list --tags --max-count=1`)" .
