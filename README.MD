# Introduction

Welcome to the server for the NVDARemote addon, written in the [Go programming language.](https://golang.org/) The idea for this program was enspired by [the one released here.](https://github.com/jmdaweb/NVDARemoteServer)

The original server for the addon is written in [Python,](https://www.python.org/) which works well enough under most circumstances. However, there were a few reasons why I wanted to write the server in Go.

- Performance. Go can be much faster than Python if you write your program properly.
- Less memory should be used by the Go program.
- Easy to compile on other operating systems without changing the code base much, if at all.
- Can be compiled into a static, position-independent binary if desired, for use across systems using different libraries.
- Should be able to easily upgrade a program to utalize the latest version of Go without changing API calls.


## Goals

My goals for this project are to create a stable server for the NVDARemote addon, which will be fast and have an efficient memory footprint. Configuration of the program will be done through the command line, which is more efficient than parsing a configuration file.


# Installation

At the moment, there is one method for installing the program. With the latest version of Go installed, use the following command.

```console
$ go get github.com/tech10/nvdaRemoteServer
```

This should download and compile the latest code, placing the binary within your GOBIN environment variable. Presuming you have the GOBIN in your path, you will be able to execute it fairly easily. If not, you will need to update your path or execute the binary with its full path.


# Usage

```console
$ nvdaRemoteServer -cert="/path/to/ssl/certificate" -key="/path/to/ssl/key" [-ip4=""] [-ip6=""] [-port=6837] [-loglevel=0]
```

Please note that the brackets around a parameter indicate that it is optional.


## Parameters

### Required

#### `-cert`

This is the path to the SSL certificate the program will use to communicate securely, as the NVDA addon uses TLS for secure communication.


#### `-key`

This is the path to the SSL key.


### Optional

#### `-ip4`

IPV4 address for the program to listen for incoming connections on. By default, all addresses are used. If you don't want to listen on IPV4, enter `-ip4="0"` as the parameter.


#### `-ip6`

IPV6 address for the program to listen for incoming connections on. By default, all addresses are used.  If you don't want to listen on IPV6, enter `-ip6="0"` as the parameter.


#### `-port`

Choose the port the server will use. At the moment, this isn't configurable by the NVDA addon itself. You should probably leave it at 6837.


#### `-loglevel`

This will choose what you want logged. The default is 0.


##### Logging Levels

Each level above the previous will also log what the prior level is logging. For instance, 1 will log both levels 0 and 1.

- -1 will disable logging.
- 0 will log when the server has started, stopped, or if a serious error has occurred. This should also log program panicks.
- 1 will log information about which clients connect, logging both their ID and IP address.
- 2 will log what channels each client joins and leaves, which will contain channel passwords. Don't use this log level in production.
- 3 will log what the program is doing at each stage of its operation. Use this for debugging purposes only.
- 4 will log the protocol that the server and client are exchanging. Don't use this unless you're a developer or you want to annalize the protocol being used. This will also cause a performance degrodation.


# Bugs

Open an issue explaining what the bug is and how you encountered it. Be detailed, or your issue will be closed if it can't be resolved properly.


# Contributing

Fork this project and submit a pull request. Please use another branch on your fork of this project if you are submitting a pull request for something. Also, keep the following guidelines in mind.

- Test your contribution before submitting them, making sure the program compiles properly.
- Remember to make use of gofmt. This will keep the formatting of the code standard for everyone.
- Try and keep your code as clean and efficient as possible.


# Final thoughts

Primarily, I am writing this program for my use, but am releasing it for anyone to utalize, should they wish. At the moment, I'm not providing any specific compiled or packaged releases. That may change in the future.