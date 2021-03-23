# Introduction

Welcome to the server for the NVDARemote addon, written in the [Go programming language.](https://golang.org/) The idea for this program was enspired by [the one released here.](https://github.com/jmdaweb/NVDARemoteServer)

The original server for the addon is written in [Python,](https://www.python.org/) which works well enough under most circumstances. However, there were a few reasons why I wanted to write the server in Go.

- Performance. Go can be much faster than Python if you write your program properly.
- Automatic generation of the self-signed certificate, if desired. Both the addon and the server utalize self-signed SSL certificates, and the addon doesn't verify a certificate's authenticity before connecting. Therefore, I wanted a program that would automatically generate the needed certificate, store it in memory, and use it until the program is terminated.
- Less memory should be used by the Go program.
- Easy to compile on other operating systems without changing the code base much, if at all.
- Can be compiled into a static, position-independent binary if desired, for use across systems using different libraries.
- Should be able to easily upgrade a program to utalize the latest version of Go without changing API calls.


## Goals

My goals for this project are to create a stable server for the NVDARemote addon, which will be fast and have an efficient memory footprint. Configuration of the program will be done through the command line. Any log output can easily be redirected by the various facilities across operating systems, so a log file wouldn't need to be implemented by the program.


# Installation

At the moment, there is one method for installing the program. With the latest version of Go installed, use the following command.

```console
$ go install github.com/tech10/nvdaRemoteServer
```

This should download and compile the latest code, placing the binary within your GOBIN environment variable. Presuming you have the GOBIN in your path, you will be able to execute it fairly easily. If not, you will need to update your path or execute the binary with its full path.


# Usage

```console
$ nvdaRemoteServer [-cert /path/to/ssl/certificate] [-key /path/to/ssl/key] [-ip4 0.0.0.0] [-ip6 [::]] [-port 6837] [-motd "Example message of the day."] [-motd-always-display false] [-loglevel 0]
```

Please note that the brackets around a parameter indicate that it is optional.


## Parameters

### Optional

#### `-cert`

This is the path to the SSL certificate the program will use to communicate securely, as the NVDA addon uses TLS for secure communication.


#### `-key`

This is the path to the SSL key. Both the certificate and key file need to exist and be accessible by the program, or the program will fall back to generating its own certificate.


##### Note about the cert and key files

If you are using the official NVDA addon as a server, or the unofficial one I linked above, you are welcome to use the same server.pem file for the certificate and key. The server should load successfully under this configuration.

If the certificate and key files both exist and fail to load a valid SSL key pair, the program will terminate rather than falling back on automatic self-signed SSL key generation.


#### `-ip4`

IPV4 address for the program to listen for incoming connections on. By default, all addresses are used. If you don't want to listen on IPV4, enter `-ip4 0` as the parameter.


#### `-ip6`

IPV6 address for the program to listen for incoming connections on. By default, all addresses are used.  If you don't want to listen on IPV6, enter `-ip6 0` as the parameter.

If you use 0 for both IPV4 and IPV6 addresses, the program will fall back to listening on all IP addresses.

The IPV6 parameter must format an IPV6 address in brackets. The address must also be on an interface for your computer. For example, this type of parameter can be used.

`[fd80::ffe8]`

So long as that is a valid IPV6 address on one of your network interfaces, this example in the local prefix of IPV6 addresses, you will be able to listen for incoming connections to the server.


#### `-port`

Choose the port the server will use. At the moment, this isn't configurable by the NVDA addon itself. You should probably leave it at 6837. Valid port numbers are between 1 and 65536.


#### `-motd`

Enter a message of the day for your server. You probably want to quote this string in the shell so spaces will be escaped properly, as in the example command line parameter.


#### `-motd-always-display`

Force the client to display the message of the day from the server, even if it hasn't changed since you last connected. This value is a boolean, so it can be true, false, 1, or 0. The numbers 1 and 0 are the same as true and false.


#### `-loglevel`

This will choose what you want logged. The default level is 0.


##### Logging Levels

Each level above the previous will also log what the prior level is logging. For instance, 1 will log both levels 0 and 1.

- -1 will disable logging, with the exception of error messages, which are always logged. This includes panicks, which crash the program.
- 0 will log when the server has started, stopped, or if an error has occurred that isn't severe enough to be logged at all times.
- 1 will log information about which clients connect, logging both their ID and IP address.
- 2 will log what channels each client joins and leaves, which will contain channel passwords. Don't use this log level in production.
- 3 will log what the program is doing at each stage of its operation. Use this for debugging purposes only.
- 4 will log the protocol that the server and client are exchanging. Don't use this unless you're a developer or you want to annalize the protocol being used. This might cause a performance degrodation, since the protocol being exchanged is also sent to the console, or redirected to a file by your operating system if you've told it to do so.


# Notes on self-signed certificate generation

The certificate that this program generates will allow for secure verification. However, like the certificate packaged by the addon, you can't verify it by using any certificate authority. If this changes in the addon and it requires certificate verification in the future, a certificate, key, and a root certificate authority file will be required for the server to use. The client would then require the appropriate files to verify the certificate is valid, unless you get a certificate from Letsencrypt, or another certificate authority.

To generate the certificate, the program will use a source of random entropy that is cryptographically secure. This can be a problem on servers that are headless, meaning you access them remotely only. However, daemons such as Haveged exist to provide your system with available entropy that can be utalized.

As a result of the secure random number generator, along with the various algorithms needed to generate the keys, the generation of a self-signed certificate can take some time, anywhere between three to ten seconds in the best case sanario. If the program takes longer than thirty seconds to generate the certificate, it is probably hanging as it waits for available entropy to complete random number generation. Once the key is generated, the program will run with its full performance.


# Bugs

Open an issue explaining what the bug is and how you encountered it. Try and be as detailed as you can, to allow the bug to be reproduced. Be detailed, or your issue will be closed if it can't be resolved properly.


# Contributing

Fork this project and submit a pull request. Please use another branch on your fork of this project if you are submitting a pull request for something. Also, keep the following guidelines in mind.

- Test your contribution before submitting them, making sure the program compiles properly.
- Remember to make use of gofmt. This will keep the formatting of the code standard for everyone.
- Try and keep your code as clean and efficient as possible.


# Final thoughts

Primarily, I am writing this program for my use, but am releasing it for anyone to utalize, should they wish. At the moment, I'm not providing any specific compiled or packaged releases. That may change in the future.
