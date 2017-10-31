#!/bin/sh

HERE=`pwd`
cd ../..
go generate ./...
cd $HERE

if which git 2>&1 > /dev/null; then
    if [ -z "`git status --porcelain`" ]; then
        STATE=clean
    else
        STATE=dirty
    fi
    GIT_VERSION=`git rev-parse HEAD`-$STATE
else
    GIT_VERSION=Unknown
fi

go build -o bootstrap_genversion
VERSION=`./bootstrap_genversion --major 1 --minor 1 --patch 1`
rm bootstrap_genversion

touch main.go
go install -v -ldflags "-X github.com/richardwilkes/gokit/cmdline.AppVersion=$VERSION -X github.com/richardwilkes/gokit/cmdline.GitVersion=$GIT_VERSION"
