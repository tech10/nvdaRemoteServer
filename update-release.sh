#!/bin/sh
# This is used by me to update the repository to a new version quickly.
# This must be executed with a new tag as the command line argument.
# If any command fails for any reason, the script will exit immediately.
# No version checks are done on the tags, and automatic commit messages are created for a version update.
# I got tired of manually updating the readme file, among other things I was doing, so combine them here.
oldtag=$(git describe --tags $(git rev-list --tags --max-count=1))
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
echo Changing readme file.
sed -i "s/$oldtag/$newtag/g" README.MD
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Generating certificate.
go run . -launch=false -log-level=-1 -gen-cert-file ./cert.pem
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Committing update.
git commit -a -s -m "release: $newtag"
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Pushing update.
git push
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Tagging update.
git tag -s -m "release: $newtag" -a $newtag
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Releasing.
goreleaser release --rm-dist
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
