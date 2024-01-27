#!/bin/bash

packages=$(go list ./... | grep -v "/mocks" | xargs)


go test $packages