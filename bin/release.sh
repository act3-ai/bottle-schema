#!/usr/bin/env bash

set -euo pipefail

help() {
    cat <<EOF
Try using one of the following commands:
prepare - prepare a release locally by producing the changelog, next version, and release notes.
approve - commit, tag, and push your approved release.
publish - publish the release to GitHub.

Dependencies: dagger, and OCI registry access.
EOF
    exit 1
}

if [ "$#" != "1" ]; then
    help
fi

set -x

case $1 in
prepare)
    if [[ $(git diff --stat) != '' ]]; then
        echo 'Git repo is dirty, aborting'
        exit 2
    fi

    dagger call release prepare export --path=.
    version=$(cat VERSION)

    echo "Please review the local changes, especially releases/$version.md"
    ;;

approve)
    # review the release material changes
    version=v$(cat VERSION)
    notesPath="releases/$version.md"
    git add VERSION CHANGELOG.md "$notesPath"
    # signed commit
    git commit -S -m"chore(release): prepare for $version"
    # annotated and signed tag
    git tag -s -a -m "Official release $version" "$version"
    # push this branch and the associated tags
    git push --follow-tags
    ;;

publish)
    # CI can then run this task (or it can be run manually)
    dagger call publish --token=env:GITHUB_TOKEN
    ;;

*)
    help
    ;;
esac
