#!/bin/bash

# --------------------------
# LINT THE CODE
# --------------------------
if [ "$1" = "lint" ]; then
    if [ "$2" = "fix" ]; then
        echo "####################################################################"
        echo "## Check and fix code format - golangci-lint"
        echo "####################################################################"
        pushd "$(dirname "$0")" > /dev/null
        cd ..

        goimports -w .
        # Auto-fix supported issues (mainly formatting and imports)
        golangci-lint run ./... --fix

        popd > /dev/null
    else
        echo "####################################################################"
        echo "## Check code format - golangci-lint"
        echo "####################################################################"
        pushd "$(dirname "$0")" > /dev/null
        cd ..

        # Just check (no fixing)
        golangci-lint run ./... 

        popd > /dev/null
    fi
else
    echo "No valid command passed"
    echo 'Use "lint" or "lint fix"'
fi
