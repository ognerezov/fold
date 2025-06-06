#!/bin/bash
args=()

# Forward all non-empty variables starting with INPUT_
for var in $(printenv | grep -o '^INPUT_[^=]*'); do
    value="${!var}"
    if [ -n "$value" ]; then
        # Remove INPUT_ prefix and convert to lowercase with hyphens
        arg_name=$(echo "${var#INPUT_}" | tr '[:upper:]' '[:lower:]' | tr '_' '-')
        args+=("--${arg_name}=${value}")
    fi
done

./fold "${args[@]}"