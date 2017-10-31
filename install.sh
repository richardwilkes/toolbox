#!/bin/sh

ROOT=`pwd`

cd $ROOT/cmdline/cmd/genversion
./install.sh

cd $ROOT/i18n/cmd/go-i18n
./install.sh
