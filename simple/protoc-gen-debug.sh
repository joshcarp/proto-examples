#!/bin/env bash

REQ=$(cat /dev/stdin | buf convert --type=google.protobuf.compiler.CodeGeneratorRequest | jq -C)
echo '{"file": [{"name":"CodeGenerateRequest.json", "content": "$REQ"}]}' | buf convert --type=google.protobuf.compiler.CodeGeneratorResponse --from -#format=json --to=-#format=bin