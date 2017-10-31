#!/bin/sh

if which genversion 2>&1 > /dev/null; then
    VERSION=`genversion --major 1`
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
    touch main.go
    go install -v -ldflags "-X github.com/richardwilkes/gokit/cmdline.AppVersion=$VERSION -X github.com/richardwilkes/gokit/cmdline.GitVersion=$GIT_VERSION"
else
    echo You must install genversion first:
    echo ""
    echo "    go get -u github.com/richardwilkes/gokit"
    echo "    cd $GOPATH/src/github.com/richardwilkes/gokit"
    echo "    ./install.sh"
fi
