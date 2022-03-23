#!/usr/bin/env bash

cargo run --quiet -- --version | awk '{print $2}'
