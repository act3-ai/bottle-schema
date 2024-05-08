#!/usr/bin/env bash

set -e

ver="$1"

gorelease -base="$(git describe --abbrev=0 --tags)" -version="v$ver"

