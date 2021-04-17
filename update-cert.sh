#!/bin/sh
# You need go installed to run this script.
# It will regenerate the certificate included with this program.
# If this takes a while, your system likely lacks entropy. See the readme for an explanation of this.
. ./functions.sh
echo Updating certificate.
gen_cert
echo Successfully updated certificate.
