#!/bin/sh
set -e

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"

resolve_webrpc_gen() {
    if [ -n "$WEBRPC_GEN_BIN" ]; then
        "$WEBRPC_GEN_BIN" "$@"
    else
        go -C "$REPO_ROOT/tools" tool webrpc-gen "$@"
    fi
}

resolve_webrpc_test() {
    if [ -n "$WEBRPC_TEST_BIN" ]; then
        "$WEBRPC_TEST_BIN" "$@"
    else
        go -C "$REPO_ROOT/tools" tool webrpc-test "$@"
    fi
}

cleanup() {
    echo "Cleaning up..."
    # Kill all child processes of this script
    pkill -P $$
}

trap cleanup EXIT

cd "$SCRIPT_DIR"

(
    resolve_webrpc_test -print-schema
) > test.ridl
(
    resolve_webrpc_gen -schema="$SCRIPT_DIR/test.ridl" -client -okhttpTransport=true -out="$SCRIPT_DIR/src/main/kotlin/TestApiClient.kt" -target="$REPO_ROOT"
)
(
    resolve_webrpc_test -server -port=9988 -timeout=10m
) &
until nc -z localhost 9988; do sleep 0.1; done; 
./gradlew --no-daemon test
