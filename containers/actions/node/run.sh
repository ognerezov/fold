#!/bin/bash
set -eo pipefail

# Process backend arguments from INPUT_* environment variables
args=()

for var in $(printenv | grep -o '^INPUT_[^=]*'); do
    # Skip special variables for frontend and test commands
    if [[ "$var" == "INPUT_RUN" || "$var" == "INPUT_TEST" || "$var" == "INPUT_START_DELAY" || "$var" == "INPUT_WORK_DIR"  ]]; then
        continue
    fi

    value="${!var}"
    if [ -n "$value" ]; then
        # Remove INPUT_ prefix and convert to lowercase with hyphens
        arg_name=$(echo "${var#INPUT_}" | tr '[:upper:]' '[:lower:]' | tr '_' '-')
        args+=("--${arg_name}=${value}")
    fi
done

# Set default commands if not provided
frontend_cmd="${INPUT_RUN}"
test_cmd="${INPUT_TEST:-npm run test}"
start_delay="${INPUT_START_DELAY:5}"
echo "$start_delay"

# Temporary files for output
backend_log=$(mktemp)
frontend_log=$(mktemp)
test_output=$(mktemp)

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    kill "$backend_pid" "$frontend_pid" 2>/dev/null || true
    rm -f "$backend_log" "$frontend_log" "$test_output"
}

# Trap signals to ensure cleanup
trap cleanup EXIT SIGINT SIGTERM

# Start backend process
backend_cmd="./fold ${args[@]}"
echo "Starting backend: $backend_cmd"
eval "$backend_cmd" > "$backend_log" 2>&1 &
backend_pid=$!

if [[ -z "${INPUT_WORK_DIR}" ]]; then
    echo "workdir not set"
else
    cd "${INPUT_WORK_DIR}"
fi

if [[ -z "${INPUT_RUN}" ]]; then
    echo "Run command not provided. not running frontend server"
else
    # Start frontend process
    echo "Starting frontend: $frontend_cmd"
    eval "$frontend_cmd" > "$frontend_log" 2>&1 &
    frontend_pid=$!
fi

# Wait for services to become available
echo "Waiting for services to start..."
sleep 5  # Adjust as needed for your services
# Run tests and capture output
echo "Running tests: $test_cmd"
set +e  # Temporarily allow errors so we can capture test exit code
eval "$test_cmd" 2>&1 | tee "$test_output"

test_exit=${PIPESTATUS[0]}
set -e


# Output logs if tests failed
if [ "$test_exit" -ne 0 ]; then
    echo -e "\n=== Backend output ===" >&2
    cat "$backend_log" >&2
    echo -e "\n=== Frontend output ===" >&2
    cat "$frontend_log" >&2
    echo -e "\n=== Test output ===" >&2
    cat "$test_output" >&2
fi

# Output result in GitHub Actions format
cat "$test_output" >> "$GITHUB_OUTPUT"

exit "$test_exit"