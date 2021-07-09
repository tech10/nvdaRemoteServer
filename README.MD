# Introduction

Welcome to the server for the NVDARemote addon, written in the [Go programming language.](https://golang.org/) The idea for this program was enspired by [the one released here.](https://github.com/jmdaweb/NVDARemoteServer)

The original server for the addon is written in [Python,](https://www.python.org/) which works well enough under most circumstances. However, there were a few reasons why I wanted to write the server in Go.

- Performance. Go can be much faster than Python if you write your program properly.
- Automatic generation of the self-signed certificate, if desired. Both the addon and the server utalize self-signed SSL certificates, and the addon originally didn't verify a certificate's authenticity before connecting. This has since been updated. I wanted a program that would automatically generate the needed certificate, store it in memory, and use it until the program is terminated.
- Less memory should be used by the Go program.
- Easy to compile on other operating systems without changing the code base much, if at all.
- Can be compiled into a static, position-independent binary if desired, for use across systems using different libraries.
- Should be able to easily upgrade a program to utalize the latest version of Go without changing API calls within the program.


## Goals

My goals for this project are to create a stable server for the NVDARemote addon, which will be fast and have an efficient memory footprint. Configuration of the program will be done through the command line. Any log output can easily be redirected by the various facilities across operating systems, so a log file wouldn't need to be implemented by the program.


# Downloading / installation

## Go

With the latest version of Go installed, use the following command.

```console
$ go install github.com/tech10/nvdaRemoteServer
```

This should download and compile the latest code, placing the binary within your GOBIN environment variable. Presuming you have the GOBIN in your path, you will be able to execute it fairly easily. If not, you will need to update your path or execute the binary with its full path.


### Static build

If you require a completely static build of the program, you can clone the GitHub repository and execute the bash script. So long as you have [musl libc](https://www.musl-libc.org/) installed, and musl-gcc in your path, you can do the following:

```console
$ git clone https://github.com/tech10/nvdaRemoteServer
$ cd ./nvdaRemoteServer
$ ./build-static.sh
```

Installing musl libc is beyond the scope of this document, and is specific to your Linux distribution.


## Docker

Thanks to a feature request in [issue 1,](https://github.com/tech10/nvdaRemoteServer/issues/1) a Dockerfile has been added. Automatically built and published images with GitHub actions are currently placed on the GitHub Container Registry, a service that is in public beta. As such, this service is subject to change until it's announced as stable. I don't anticipate many changes to the service, so everything ought to continue working as is.

All example commands expose the host networking to the docker image, which is the quickest means of hosting the service without any difficulty. There are other ways, but they won't work as well with IPV6 unless you set it up, and are beyond the scope of this document.

By default, the docker image will contain the certificate within the GitHub repository, and utalize it. This will prevent entropy difficulties on systems that may not be able to generate their own certificates quickly. This can be altered by the user at the run time of the Docker image. Presuming your image is tagged nvdaremoteserver, here is a sample command.

```console
$ docker run --network host nvdaremoteserver /nvdaRemoteServer
```

This will run the Docker image with no parameters, allowing the server to choose its own defaults. This will prevent it from automatically using the included certificate within the git repository. To use the included certificate, you can use the minimal examples below.


### Downloading the Docker image

#### GitHub Container Registry

The automatically built images are [located here.](https://github.com/users/tech10/packages/container/package/nvdaremoteserver-docker) Provided on that page is a command you can copy to the clipboard which should pull the latest image. Here is an example use of downloading and running the image.

```console
$ docker pull ghcr.io/tech10/nvdaremoteserver-docker:v1.0.0
$ docker run --network host ghcr.io/tech10/nvdaremoteserver-docker:v1.0.0
```


#### Docker Hub

Images are available on [Docker Hub.](https://hub.docker.com) They are not automatically built as they once were, due to that feature leaving free accounts on June 18, 2021. The images may not be as up to date as those on GitHub. Here is an example to pull the latest image and run it.

```console
$ docker pull tech10/nvdaremoteserver
$ docker run --network host tech10/nvdaremoteserver
```


### Manually build docker image

Clone the repository, and from within the directory, build a Docker image, then run it. Some sample commands are below, one of which will update the certificate from the repository, rewriting it to a freshly generated certificate. To update the certificate, you need go installed, as documented within the shell script. If you simply want to build the docker image, remove the command updating the certificate. You only need Docker installed in order to build the image.

```console
$ git clone https://github.com/tech10/nvdaRemoteServer
$ cd ./nvdaRemoteServer
$ ./update-cert.sh
$ docker build -t nvdaremoteserver-docker .
$ docker run --network host nvdaremoteserver-docker
```


# About the included certificate

The included certificate file is a single file that contains an automatically generated self-signed certificate from this program, along with its private key. They are both encoded in the same pem format that the official addon uses. Before every release, a new certificate file is generated and uploaded to the GitHub repository. For now, this is how it is placed in the Docker images, and how it is made available to other users. This is subject to change in the future.


# Usage

```console
$ nvdaRemoteServer [-pid-file /path/to/pid/file] [-conf-file /path/to/configuration/file] [-conf-read=true] [-gen-conf-file /path/to/generated/configuration/file] [-gen-conf-dir=false] [-create=false] [-address :6837] [-cert-file /path/to/ssl/certificate] [-key-file /path/to/ssl/key] [-gen-cert-file /path/to/created/cert/file] [-motd "Example message of the day."] [-motd-always-display=false] [-send-origin=true] [-log-level=0] [-log-file /path/to/log/file] [-launch=true]
```

Please note that the brackets around a parameter indicate that it is optional.


## Parameters

### Optional

#### `-conf-file`

This is a path to an existing configuration file. If the configuration file can't be read, the program will alert you and exit with an error.

When reading a configuration file, all command line parameters take priority over anything within a configuration file. For example, if you create a configuration file, then later decide you wish to listen on a different address, the address you specify, presuming it isn't the default address, will be used over that in the configuration file.

Configuration files are searched for automatically in two places if this parameter is not supplied. First, if a file named nvdaRemoteServer.json is found in the current working directory, it will be read. Second, the users configuration directory will be searched for a directory named nvdaRemoteServer, and a configuration file with the previously stated name. If neither of these files are found, and you don't specify a configuration file, the program will continue execution.

If a configuration file you specify is invalid, the program will exit after telling you what error has been encountered. If you haven't specified a configuration file, but one is found in one of the searched directories that is invalid, you will be alerted and the program will continue execution with any given command line parameters.

If a configuration file is successfully read, and isn't in the current working directory, the program will change its working directory to that of the configuration file that was read. Any relative parameters to files, such as nvdaRemoteServer.log for a log file, will be created and written to under the configuration file directory.

If a configuration file is read, and using default parameters, you will be alerted that the configuration file is using default parameters and none is needed. The program will then continue execution.


#### `-conf-read`

This will choose whether or not the program will read a configuration file. If set to false, no configuration file will be read, the program will warn you of this, then continue execution. If you have set the `-conf-file` parameter, it will be reset to its default value and the program will warn you of its reset.


#### `-gen-conf-file`

This is a path to a configuration file the program will attempt to generate from given command line parameters. If you have only specified the generation of a configuration file, no configuration file will be generated, you will be alerted on the info log level, and the program will continue execution with the default parameters.

When generating a configuration file, you can automatically specify a generated certificate file that will be used for the cert and key parameters. This can be done by specifying the `-gen-cert-file` parameter, but not specifying the `-cert-file` or `-key-file` parameters.

If the configuration file generation is successful, the working directory will be changed to that of the configuration file, if different than the current working directory. This will occurr if creating a user configuration, for example.


#### `-gen-conf-dir`

This will generate a configuration directory for the currently running user if set to true. If the directory doesn't exist, it will be automatically created, if possible. If it can't be created, you will be alerted with an error and the program will abort execution. The program will tell you what path is used for the generation of the configuration file and directory, so you can find it later if you desire.


#### `-create`

If set to true, this parameter will attempt to create directories upon any operation that requires writing to a file. The default is false.

This parameter is temporarily set to true if you specify that a user configuration directory is to be generated.


#### `-pid-file`

Path to a file where the process ID is stored.

This file will only be created once the server has successfully started, and will be removed upon shutdown. This could be useful if you've started the program in the background, and wish to kill the process without searching for its process ID by name. If the program fails to create the file, it will warn you via the debug log level and continue execution.


#### `-address`

Address for the program to listen for incoming connections on in the form ip:port. By default, all addresses are used, and the server accepts connections on port 6837. This can be declared more than once.

If you want to listen on an IPV6 address, the address must be surrounded by brackets. The address must also be on an interface for your computer. For example, this type of parameter can be used.

`[fd80::ffe8]:6837`

So long as that is a valid IPV6 address on one of your network interfaces, this example in the local prefix of IPV6 addresses, you will be able to listen for incoming connections to the server. To listen on only IPV6 addresses, use the following example.

`[::]:6837`

To listen on an IPV4 address, it's as simple as using a valid parameter such as the following example, which will listen on all IPV4 addresses only.

`0.0.0.0:6837`

The default port, if no address is declared, is 6837. Valid port numbers are between 1 and 65536. When declaring an address for the server to listen on, you must also declare a port, or the parameter will be invalid. You need not declare an address, however. For example, if you want to listen on all addresses, but use port 5000, use the following.

`:5000`


#### `-cert-file`

This is the path to the SSL certificate the program will use to communicate securely, as the NVDA addon uses TLS for secure communication.


#### `-key-file`

This is the path to the SSL key. Both the certificate and key file need to exist and be accessible by the program, or the program will fall back to generating its own certificate.


##### Note about the cert and key files

If you are using the official NVDA addon as a server, or the unofficial one I linked above, you are welcome to use the same server.pem file for the certificate and key. The server should load successfully under this configuration. If you choose to generate a certificate file, the key and certificate will be in a single file, just as they are with the official addon.

If the certificate and key files both exist and fail to load a valid SSL key pair, the program will terminate rather than falling back on automatic self-signed SSL key generation.


#### `-gen-cert-file`

Path to a location where a file can be written with the automatically generated certificate and key.

When the program generates its own self-signed certificate, you can optionally write this certificate to a file, so as to easily use it again in the future. The certificate and key will be written to a single file. If the file can't be written, the program will warn you via the debug log level and continue execution.


##### Notes on self-signed certificate generation

The certificate that this program generates will allow for secure verification. However, like the certificate packaged by the addon, you can't verify it by using any certificate authority. If you know that you have generated the certificate, you can allow the addon to connect by trusting its fingerprint. Alternatively, you could get a verified certificate from Letsencrypt, or another certificate authority and use that, as the addon will make sure the certificate can be verified as secure.

To generate the certificate, the program will use a source of random entropy that is cryptographically secure. This can be a problem on servers that are headless, meaning you access them remotely only. However, daemons such as Haveged exist to provide your system with available entropy that can be utalized.

As a result of the secure random number generator, along with the various algorithms needed to generate the keys, the generation of a self-signed certificate can take some time, anywhere between three to ten seconds in the best case sanario. If the program takes longer than thirty seconds to generate the certificate, it is probably hanging as it waits for available entropy to complete random number generation. Once the key is generated, the program will run with its full performance.

If the self-signed certificate takes a particularly long time to create on your system, you can elect to write the certificate to a file as documented above, and optionally do only that by setting the launch parameter to false, which is documented below. You could then launch the server with the generated certificate file as the parameter for the cert and key parameters, preventing long start times on systems with slow processors or entropy difficulties. Alternatively, you could have someone else generate a certificate for you with this program, then send it to you via a secure method of transfer. The final method would be using the certificate file included in this repository, though this may be less secure than the other methods outlined. This is no more secure than the official addon, which packages its own certificate rather than generating it.


#### `-motd`

Enter a message of the day for your server. You probably want to quote this string in the shell, ensuring spaces will be escaped properly, as in the example command line parameter.


#### `-motd-always-display`

Force the client to display the message of the day from the server, even if it hasn't changed since you last connected. This value is a boolean, so it can be true, false, 1, or 0. The numbers 1 and 0 are the same as true and false.


#### `-send-origin`

By default, when the server receives a message from a client, it will send that same message to all clients that need to receive it, but it will add an origin field to that message. This requires that the message be decoded into a value the program can more easily manipulate, an origin message added to it, then encoded back to the value that will be sent to all clients. This could cause a slight performance hit. You can disable this feature by setting it to false, if desired, though you might find some things don't work properly for you if you do so. If you set this to false, the server will warn yu that it may impact the functionality of clients when the origin field is required.


#### `-log-file`

Choose a file for the program to log its data. Any logged information will always be sent to the console, but in addition to this, a log file can also be used. By default, logged data is only sent to the console.

If the server is not being launched, no log file will be written, the program will warn you, then continue execution.


#### `-log-level`

This will choose what you want logged. The default level is 0.


##### Logging Levels

Each level above the previous will also log what the prior level is logging. For instance, 1 will log both levels 0 and 1.

- -1 will disable logging, with the exception of error messages, which are always logged. This includes panicks, which crash the program.
- 0 will log when the server has started, stopped, or if an error has occurred that isn't severe enough to be logged at all times.
- 1 will log information about which clients connect, logging both their ID and IP address.
- 2 will log what channels each client joins and leaves, which will contain channel passwords. Don't use this log level in production.
- 3 will log what the program is doing at each stage of its operation. Use this for debugging purposes only.
- 4 will log the protocol that the server and client are exchanging. Don't use this unless you're a developer or you want to annalize the protocol being used. This might cause a performance degrodation, since the protocol being exchanged is also sent to the console, or redirected to a file by your operating system if you've told it to do so.


#### `-launch`

By default, the server will attempt to launch itself if nothing has caused it to abort execution prematurely. If you set this parameter to false, the server will shut down immediately after most configuration is complete. This can be useful if you only wanted to generate a self-signed certificate and write it to a file, for instance.


## Other parameters

### `version`

Print the program version, then shut down immediately. Example:

```console
$ nvdaRemoteServer version
development
$ 
```


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


## Observations using the servers

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

- Test your contributions before submitting them, making sure the program compiles properly.
- Remember to make use of gofmt. This will keep the formatting of the code standard for everyone.
- Try and keep your code as clean and efficient as possible.


# Final thoughts

Primarily, I am writing this program for my use, but am releasing it for anyone to utalize, should they wish.
