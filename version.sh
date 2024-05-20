#!/bin/sh -e
#
# Self-modifying script that updates the version numbers
# Requires "setconf"
#

# The current version goes here, as the default value
VERSION=${1:-'1.0.10'}

if [ -z "$1" ]; then
  echo "The current version is $VERSION, pass the new version as the first argument if you wish to change it"
  exit 0
fi

echo "Setting the version to $VERSION"

# Set the version in various files
setconf README.md '* Version' $VERSION
setconf main.go versionString "\"Desktop File Generator "$VERSION"\""

# Update the date and version in the man page
d=$(LC_ALL=C date +'%d %b %Y')
sed -i "s/\"[0-9]* [A-Z][a-z]* [0-9]*\"/\"$d\"/g" gendesk.1
sed -i "s/[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+/$VERSION/g" gendesk.1

# Update the version in this script
sed -i "s/[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+/$VERSION/g" "$0"

echo 'Also update the changelog in README.md'
