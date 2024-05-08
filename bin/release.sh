#!/usr/bin/env bash

set -e

ver="$1"
export ver

echo "$ver" > VERSION
