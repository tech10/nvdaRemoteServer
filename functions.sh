#!/bin/sh
check() {
eval $@
if [ $? -ne 0 ]; then
echo "Command $1 failed to execute."
exit 10
fi
}
git_version() {
git describe --tags $(git rev-list --tags --max-count=1)
}
