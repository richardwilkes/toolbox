#! /usr/bin/env bash

set -eo pipefail

trap 'echo -e "\033[33;5mBuild failed on build.sh:$LINENO\033[0m"' ERR

RACE=-race

# Process args
for arg in "$@"
do
  case "$arg" in
    --skip-linters|-l)    SKIP_LINTERS=1 ;;
    --skip-test|-t)       SKIP_TESTS=1 ;;
    --omit-race|-r)       RACE= ;;
    --fast|-f)            SKIP_LINTERS=1; SKIP_TESTS=1 ;;
    --help|-h)
      echo "$0 [options]"
      echo "  -f, --fast             Same as -l -t"
      echo "  -l, --skip-linters     Skip linters"
      echo "  -t, --skip-tests       Skip tests"
      echo "  -r, --omit-race        Omit the -race option in tests"
      echo "  -h, --help             This help text"
      exit 0
      ;;
    *) echo "Invalid argument: $arg"; BAIL=1 ;;
  esac
done
if [ -n "$BAIL" ]; then
  exit 1
fi

# Setup the tools we'll need
TOOLS_DIR=$PWD/tools
mkdir -p "$TOOLS_DIR"
if [ -z $SKIP_LINTERS ]; then
  if [ ! -e "$TOOLS_DIR/golangci-lint" ] || [ "$("$TOOLS_DIR/golangci-lint" version 2>&1 | awk '{ print $4 }' || true)x" != "1.22.2x" ]; then
    echo -e "\033[33mInstalling version 1.22.2 of golangci-lint into $TOOLS_DIR...\033[0m"
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$TOOLS_DIR" v1.22.2
  fi
fi
export PATH=$TOOLS_DIR:$PATH

# Setup version info
if command -v git 2>&1 > /dev/null; then
  if [ -z "$(git status --porcelain)" ]; then
    STATE=clean
  else
    STATE=dirty
  fi
  GIT_VERSION=$(git rev-parse HEAD)-$STATE
  GIT_TAG=$(git tag --points-at HEAD)
  if [ -z $GIT_TAG ]; then
    GIT_TAG=$(git tag --list --sort -version:refname | head -1)
    if [ -n "$GIT_TAG" ]; then
      GIT_TAG=$GIT_TAG~
    fi
  fi
  if [ -n "$GIT_TAG" ]; then
    VERSION=$(echo $GIT_TAG | sed -E "s/^v//")
  else
    VERSION=""
  fi
else
  GIT_VERSION=Unknown
  VERSION=""
fi

# Build the code
echo -e "\033[33mBuilding Go code...\033[0m"
LINK_FLAGS="-X github.com/richardwilkes/toolbox/cmdline.AppVersion=$VERSION"
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.BuildNumber=$(date "+%Y%m%d%H%M%S")"
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.GitVersion=$GIT_VERSION"
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.CopyrightYears=2016-$(date "+%Y")"
find . -iname "*_gen.go" -exec /bin/rm {} \;
go generate ./...
go build -v -ldflags=all="$LINK_FLAGS" ./...

# Run the linters
if [ -z $SKIP_LINTERS ]; then
  echo -e "\033[33mRunning Go linters...\033[0m"
  golangci-lint run
else
  echo -e "\033[33mSkipping Go linters\033[0m"
fi

# Run the tests
if [ -z $SKIP_TESTS ]; then
  echo -e "\033[33mRunning tests...\033[0m"
  go test $RACE ./...
else
  echo -e "\033[33mSkipping tests\033[0m"
fi

# Install executables
echo -e "\033[33mInstalling executables...\033[0m"
go install -ldflags=all="$LINK_FLAGS" ./xio/fs/mkembeddedfs
go install -ldflags=all="$LINK_FLAGS" ./i18n/i18n
