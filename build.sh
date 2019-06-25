#!/bin/sh

# This will install the various tools provided with this repo. The only
# difference between running this and doing:
#
# go get -u github.com/richardwilkes/toolbox/...
#
# is proper version numbers with build dates and git revisions will be
# embedded into the resulting executables.

ROOT=`pwd`

if which git 2>&1 > /dev/null; then
    if [ -z "`git status --porcelain`" ]; then
        STATE=clean
    else
        STATE=dirty
    fi
    GIT_VERSION=`git rev-parse HEAD`-$STATE
    GIT_TAG=`git tag --points-at HEAD`
    MAJOR=`echo $GIT_TAG | sed -E "s/v([0-9]+)\\.([0-9]+)\\.([0-9]+)/\1/"`
    MINOR=`echo $GIT_TAG | sed -E "s/v([0-9]+)\\.([0-9]+)\\.([0-9]+)/\2/"`
    PATCH=`echo $GIT_TAG | sed -E "s/v([0-9]+)\\.([0-9]+)\\.([0-9]+)/\3/"`
else
    GIT_VERSION=Unknown
fi
if [ -z $MAJOR ]; then
    MAJOR=0
fi
if [ -z $MINOR ]; then
    MINOR=0
fi
if [ -z $PATCH ]; then
    PATCH=0
fi

go generate ./...

VERSION=`go run cmd/genversion/main.go --major $MAJOR --minor $MINOR --patch $PATCH`
LINK_FLAGS="-X github.com/richardwilkes/toolbox/cmdline.AppVersion=$VERSION"
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.GitVersion=$GIT_VERSION"
go install -v -ldflags=all="$LINK_FLAGS" ./...
