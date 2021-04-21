#!/bin/sh
# Generate a sample configuration file.
go run . -launch=false -gen-conf-file nvdaRemoteServer.json -cert cert.pem -key cert.pem -log-file nvdaRemoteServer.log
