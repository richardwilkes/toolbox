#! /usr/bin/env bash
set -eo pipefail

trap 'echo -e "\033[33;5mBuild failed on build.sh:$LINENO\033[0m"' ERR

GOLANGCI_LINT_VERSION=1.42.1

for arg in "$@"
do
  case "$arg" in
    --all|-a) LINT=1; TEST=1; RACE=-race ;;
    --lint|-l) LINT=1 ;;
    --race|-r) TEST=1; RACE=-race ;;
    --test|-t) TEST=1 ;;
    --help|-h)
      echo "$0 [options]"
      echo "  -a, --all  Equivalent to --lint --race"
      echo "  -l, --lint Run the linters"
      echo "  -r, --race Run the tests with race-checking enabled"
      echo "  -t, --test Run the tests"
      echo "  -h, --help This help text"
      exit 0
      ;;
    *)
      echo "Invalid argument: $arg"
      exit 1
      ;;
  esac
done

# Setup version info
if command -v git 2>&1 > /dev/null; then
  if [ -z "$(git status --porcelain)" ]; then
    STATE=clean
  else
    STATE=dirty
  fi
  GIT_VERSION=$(git rev-parse HEAD)-$STATE
  GIT_TAG=$(git tag --points-at HEAD)
  if [ -z "$GIT_TAG" ]; then
    GIT_TAG=$(git tag --list --sort -version:refname | head -1)
    if [ -n "$GIT_TAG" ]; then
      GIT_TAG=$GIT_TAG~
    fi
  fi
  if [ -n "$GIT_TAG" ]; then
    VERSION=$(echo "$GIT_TAG" | sed -E "s/^v//")
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
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.BuildNumber=$(date -u "+%Y%m%d%H%M%S")"
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.GitVersion=$GIT_VERSION"
LINK_FLAGS="$LINK_FLAGS -X github.com/richardwilkes/toolbox/cmdline.CopyrightYears=2016-$(date "+%Y")"
find . -iname "*_gen.go" -exec /bin/rm {} \;
go generate ./gen
go build -v -ldflags=all="$LINK_FLAGS" ./...

# Run the tests
if [ "$TEST"x == "1x" ]; then
  if [ -n "$RACE" ]; then
    echo -e "\033[33mTesting with -race enabled...\033[0m"
  else
    echo -e "\033[33mTesting...\033[0m"
  fi
  go test $RACE ./...
fi

# Run the linters
if [ "$LINT"x == "1x" ]; then
  TOOLS_DIR=$PWD/tools
  if [ ! -e "$TOOLS_DIR/golangci-lint" ] || [ "$("$TOOLS_DIR/golangci-lint" version 2>&1 | awk '{ print $4 }' || true)x" != "${GOLANGCI_LINT_VERSION}x" ]; then
    echo -e "\033[33mInstalling version $GOLANGCI_LINT_VERSION of golangci-lint into $TOOLS_DIR...\033[0m"
    mkdir -p "$TOOLS_DIR"
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$TOOLS_DIR" v$GOLANGCI_LINT_VERSION
  fi
  echo -e "\033[33mRunning Go linters...\033[0m"
  "$TOOLS_DIR/golangci-lint" run
fi

# Install executables
echo -e "\033[33mInstalling executables...\033[0m"
go install -ldflags=all="$LINK_FLAGS" ./i18n/i18n
