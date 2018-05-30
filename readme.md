# A zip installer for Go

## Purpose

Use this command line to download and install go packages when you don't have the ability to clone from public git repos like Github.

## How It Works

go-get-zip finds and downloads the zipped up version of the source code from GitHub, uzips the files into the `GOPATH` and installs the packages.

## Installation

In order to get started with go-get-zip clone this directory, or download and unzip the zipfile of this directory if you can't clone it. cd into the root of this directory and run `go build .`. Add the path to the executable file to your system PATH variable.

## Usage

Run `go-get-zip <package you want to install>` from the command line. For example `go-get-zip github.com/conradpcesa/go-get-zip` or `go-get-zip gopkg.in/yaml.v1`. You can also install all the dependencies in a prject by going to the directory and running `go-get-zip` (similar to running `go get` to install all go dependencies).