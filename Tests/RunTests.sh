#!/bin/sh
set -e

cleanup() {
    echo "Cleaning up..."
    # Kill all child processes of this script
    pkill -P $$
}

trap cleanup EXIT

webrpc-test -print-schema > test.ridl
webrpc-gen -schema=test.ridl -client -out=src/main/kotlin/TestApiClient.kt -target=../
webrpc-test -server -port=9988 -timeout=10m &
until nc -z localhost 9988; do sleep 0.1; done; 
./gradlew test