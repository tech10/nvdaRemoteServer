#!/bin/sh
# This is used by me to update the repository to a new version quickly.
# This must be executed with a new tag as the command line argument.
# If any command fails for any reason, the script will exit immediately.
# No version checks are done on the tags, and automatic commit messages are created for a version update.
# I got tired of manually updating the readme file, among other things I was doing, so combine them here.
. ./functions.sh
oldtag=$(git_version)
newtag=$1
if [ -z "$newtag" ]; then
echo The update tag cannot be blank.
exit
fi
if [ "$oldtag" = "$newtag" ]; then
echo The old tag and new tag are identical. Not upgrading.
exit
fi
echo "Upgrading from $oldtag to $newtag."
echo Generating certificate.
gen_cert
echo Committing update.
check git commit -a -s -m \"release: ${newtag}\"
echo Pushing update.
check git push
echo Tagging update.
check git tag -s -m \"release: ${newtag}\" -a ${newtag}
echo Releasing.
check git push origin ${newtag}
