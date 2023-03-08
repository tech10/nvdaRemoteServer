package main

import "runtime"

func buildInfo() string {
	return "This application was compiled with " + runtime.Version() + ". It was compiled for the " + runtime.GOARCH + " architecture and the " + runtime.GOOS + " operating system."
}
