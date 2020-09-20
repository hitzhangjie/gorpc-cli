#!/bin/bash -e

# step-1: format the code
gofmt -s -w .
goimports -w -local github.com .