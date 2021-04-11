#!/bin/sh
# Example script to upgrade an installed binary.
# This assumes you have tar on your system, installed the binary into /usr/bin, have sudo rights,
# and are using systemd for your daemon manager. This will automatically download the latest release.
# Modify the appropriate variables for your system, and commands, if needed.
binary="nvdaRemoteServer"
os="Linux"
arch="x86_64"
program="${binary}_${os}_${arch}"
echo Checking version.
version=$(${binary} version)
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Current version ${version}
echo Downloading.
wget -q https://github.com/tech10/nvdaRemoteServer/releases/latest/download/${program}.tar.gz
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Extracting.
tar -axf ${program}.tar.gz
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Changing ownership of binary.
sudo chown root:root ${program}/${binary}
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Moving binary to /usr/bin
sudo mv ${program}/${binary} /usr/bin/
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Restarting server.
sudo systemctl restart nvdaRemoteServer
if [ $? -ne 0 ]; then
echo Failure.
exit
fi
echo Cleaning up files.
rm -r ${program} ${program}.tar.gz
