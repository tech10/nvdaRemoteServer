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


# Statistics

These by no means should be taken as representative of any definitive stats, but the results given by systemd's tracking of memory and CPU use can speak for themselves.

When running the servers, both the Go and Python versions, I was performing similar tasks over different periods of time, typing, reading, and doing virtually everything on the remote system, including writing this section of the document and collecting the statistics.

Here is what they gathered for both the Python and Go versions of the NVDARemote server. The Python version I was running was 3.9.2, and Go was 1.16.2.


## Python

The run time was approximately one hour.

```console
$ sudo systemctl status NVDARemoteServer
 NVDARemoteServer.service - NVDARemote relay server
 Loaded: loaded (/usr/lib/systemd/system/NVDARemoteServer.service; disabled; vendor preset: disabled)
 Active: active (running) since Tue 2021-03-23 12:37:24 MDT; 1h 0min ago
 Process: 11551 ExecStart=/usr/bin/python /usr/share/NVDARemoteServer/server.py start (code=exited, status=0/SUCCESS)
 Main PID: 11553 (python)
 IP: 4.5M in, 5.2M out
 Tasks: 5 (limit: 1151)
 Memory: 12.0M
 CPU: 1min 13.385s
```


## Go

The runtime was approximately five hours, fourteen minutes.

```console
$ sudo systemctl status nvdaRemoteServer
 nvdaRemoteServer.service - NVDARemote relay server
 Loaded: loaded (/etc/systemd/system/nvdaRemoteServer.service; enabled; vendor preset: disabled)
 Active: active (running) since Tue 2021-03-23 07:21:02 MDT; 5h 14min ago
 Main PID: 7512 (nvdaRemoteServe)
 IP: 10.6M in, 11.6M out
 Tasks: 8 (limit: 1151)
 Memory: 4.8M
 CPU: 13.324s
```


## Notes about the results

There is one difference between the Python and Go versions of the server that is significant, other than the programming language being used. The Python server forks itself into the background, something that isn't strictly necessary to do with systemd processes. The Go version of the server does no forking, so the systemd service is capable of monitoring its process directly. This may cause results to be different than they should be, but the use of memory and CPU time should be fairly accurate.


## Observations of using the servers

My personal observations are the following, running both servers on a server approximately 50MS ping time away from both locations, which would make the round trip approximately 100MS:

- When running the Python version of the NVDARemote server, the delay is noticeable between the controlling computer, and the computer being controlled.
- When running the Python version of the server using [PyPy,](https://www.pypy.org/) which is a faster version of Python for longer running programs, the delay is less, better than Python and a bit more stable. I can still tell that I'm controlling a remote computer.
- When running the Go version of the NVDARemote server, the delay can still be noticed, but is less than the Python version of the server. General stability and delay between keystrokes is also improved on the Go version of the server, and sometimes, I forget that I'm actually controlling a remote computer.

I took no benchmarks of response times between sending and receiving data, but I would estimate that the Go program is at least four or five times faster than the Python program, perhaps more so. It definitely seems to use less CPU.

In comparing the two servers, keep the following in mind. Python is an interpreted language. Therefore, the program is compiled into machine code as it is executed. This compiling and reading of the program will increase CPU use and slow down the responsiveness of a program. PyPy will compile the entire Python program into machine code before it begins to execute, which makes it faster than Python for long running processes, though slightly slower in starting. I believe it was stated that PyPy is at least four times faster than Python. Go will compile the entire program into machine code before you execute it, leaving you with a binary that you will run on your computer, similar to programming languages such as C. In [one particular use case,](https://getstream.io/blog/switched-python-go/#:~:text=Go%20is%20extremely%20fast.,40%20times%20faster%20than%20Python.) it was stated that Go was forty times faster than Python.


# Bugs

Open an issue explaining what the bug is and how you encountered it. Try and be as detailed as you can, to allow the bug to be reproduced. Be detailed, or your issue will be closed if it can't be resolved properly.


# Contributing

Fork this project and submit a pull request. Please use another branch on your fork of this project if you are submitting a pull request for something. Also, keep the following guidelines in mind.

- Test your contribution before submitting them, making sure the program compiles properly.
- Remember to make use of gofmt. This will keep the formatting of the code standard for everyone.
- Try and keep your code as clean and efficient as possible.


# Final thoughts

Primarily, I am writing this program for my use, but am releasing it for anyone to utalize, should they wish. At the moment, I'm not providing any specific compiled or packaged releases. That may change in the future.
