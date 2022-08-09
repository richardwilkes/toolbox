#! /usr/bin/env bash
# Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, version 2.0. If a copy of the MPL was not distributed with
# this file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
# This Source Code Form is "Incompatible With Secondary Licenses", as
# defined by the Mozilla Public License, version 2.0.

set -eo pipefail

trap 'echo -e "\033[33;5mBuild failed on build.sh:$LINENO\033[0m"' ERR

for arg in "$@"
do
  case "$arg" in
    --all|-a) LINT=1; TEST=1; RACE=-race ;;
    --lint|-l) LINT=1 ;;
    --race|-r) TEST=1; RACE=-race ;;
    --test|-t) TEST=1 ;;
    --help|-h)
      echo "$0 [options]"
      echo "  -a, --all  Equivalent to --lint --test --race"
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

# Build the code
echo -e "\033[33mBuilding Go code...\033[0m"
go build -v ./...

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
  GOLANGCI_LINT_VERSION=1.48.0
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
go install -v ./i18n/i18n
