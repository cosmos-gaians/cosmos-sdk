#!/bin/bash

rm -r build || true
mkdir build

cargo build --release --target=wasm32-unknown-unknown

find ./target/wasm32-unknown-unknown/release/ -name "*.wasm" -exec cp "{}" build/ ";"
