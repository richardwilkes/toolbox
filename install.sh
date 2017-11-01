#!/bin/sh

# This will install the various tools provided with this repo. The only
# difference between running this and doing:
#
# go get -u github.com/richardwilkes/gokit/...
#
# is proper version numbers with build dates and git revisions will be
# embedded into the resulting executables.

ROOT=`pwd`

cd $ROOT/cmdline/cmd/genversion
./install.sh

cd $ROOT/i18n/cmd/go-i18n
./install.sh
