#!/bin/bash

LOCAL_PREFIXES=("cosmossdk.io" "github.com/cosmos/cosmos-sdk" "github.com/tabilabs/tabi")

for file in $(find . -name '*.go' -type f -path "./x/*" -not -path "./x/evm" -not -path "./x/feemarket" -not -name '*.pb.go' -not -name '*.pb.gw.go'); do
    for prefix in "${LOCAL_PREFIXES[@]}"; do
        goimports --local "$prefix" -w "$file"
    done
    gofumpt -w $file
done