#!/bin/bash
# Build
#
#
# Builds the software in 'release' mode and copies it to some location
# Depending on said location you need to run this script with 'sudo', because
# it won't check for them.

if [ -z $GOPATH ]; then
	echo "Please set GOPATH"
	exit 1
fi

if [ -z $1 ]; then
	echo Please specify destination
	exit 1
fi

go install
mv "$GOPATH/bin/shissue" "$1"
echo "Ok"
