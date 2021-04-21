#!/bin/sh
# Example script to upgrade an installed binary.
# This assumes you have tar on your system, installed the binary into /usr/bin, have sudo rights,
# and are using systemd for your daemon manager. This will automatically download the latest release.
# Modify the appropriate variables for your system, and commands, if needed.
. ./functions.sh
binary="nvdaRemoteServer"
os="Linux"
arch="x86_64"
program="${binary}_${os}_${arch}"
echo Checking version.
version=$(check ${binary} version)
echo Current version ${version}
echo Downloading.
check wget -q https://github.com/tech10/nvdaRemoteServer/releases/latest/download/${program}.tar.gz
echo Extracting.
check tar -axf ${program}.tar.gz
echo Checking new version.
new_version=$(check ${program}/${binary} version)
update() {
echo Changing ownership of binary.
check sudo chown root:root ${program}/${binary}
echo Moving binary to /usr/bin
check sudo mv ${program}/${binary} /usr/bin/
echo Restarting server.
check sudo systemctl restart nvdaRemoteServer
}
clean() {
echo Cleaning up files.
check rm -r ${program} ${program}.tar.gz
}
if [ "$version" = "$new_version" ]; then
echo The two versions are identical. Not upgrading.
clean
exit
fi
echo Upgrading from $version to $new_version
update
clean
