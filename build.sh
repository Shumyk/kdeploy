#!/bin/bash

# Build the binary and plase it in /usr/local/bin to make it available in the PATH
go build -x -o /usr/local/bin/kdeploy main.go