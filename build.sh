#!/bin/bash
cd $(dirname "${BASH_SOURCE[0]}")
OD="$(pwd)"
# Pushes application version into the build information.
_VERSION=1.0.0

build(){
	echo Packaging $1 Build
	bdir=scanner-${_VERSION}-$(date +%y%m%d)

	if [ "$4" == "0" ]; then
	    echo "build dir is "$bdir
	    rm -rf builds/$bdir && mkdir -p builds/$bdir
	fi

	GOOS=$2 GOARCH=$3 ./build.sh

	if [ "$2" == "windows" ]; then
		mv scanner builds/$bdir/scanner-$1.exe
	else
		mv scanner builds/$bdir/scanner-$1
	fi

	cp app.yaml builds/$bdir
	cp dict.yaml builds/$bdir
	cp mscsv.yaml builds/$bdir
	cp LICENSE builds/$bdir
	cp README.md builds/$bdir

	mkdir -p builds/$bdir/exiftool
    cp -r exiftool builds/$bdir
}

if [ "$1" == "all" ]; then
	rm -rf builds/
	build "windows" "windows" "amd64" 0
	build "mac" "darwin" "amd64" 1
	build "linux" "linux" "amd64" 2
	exit
fi

CGO_ENABLED=0 go build -o "$OD/scanner" main.go

#read -p "Press key to continue..." -n1 -s