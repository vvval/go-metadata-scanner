#!/bin/bash
cd $(dirname "${BASH_SOURCE[0]}")
OD="$(pwd)"
# Pushes application version into the build information.
_VERSION=1.0.0

build(){
	echo Packaging $1 Build
	bdir=scanner-${_VERSION}-$2-$3
	rm -rf builds/$bdir && mkdir -p builds/$bdir
	GOOS=$2 GOARCH=$3 ./build.sh

	if [ "$2" == "windows" ]; then
		mv scanner builds/$bdir/scanner.exe
	else
		mv scanner builds/$bdir
	fi

	cp app.yaml builds/$bdir
	cp dict.yaml builds/$bdir
	cp mscsv.yaml builds/$bdir
	cp exiftool builds/$bdir
	cp LICENSE builds/$bdir
	cp README.md builds/$bdir
}

if [ "$1" == "all" ]; then
	rm -rf builds/
	build "Windows" "windows" "amd64"
	build "Mac" "darwin" "amd64"
	build "Linux" "linux" "amd64"
	build "FreeBSD" "freebsd" "amd64"
	exit
fi

CGO_ENABLED=0 go build -o "$OD/scanner" main.go

#read -p "Press key to continue..." -n1 -s