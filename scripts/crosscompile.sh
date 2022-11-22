#!/bin/bash
version="0.1-001"
git_root_path=`git rev-parse --show-toplevel`
execution_file=curl-go
cd ${git_root_path}/scripts
#for os in linux freebsd netbsd openbsd aix android illumos ios solaris plan9 darwin dragonfly windows;
for os in linux;
do
	for arch in "amd64" "386" "arm" "arm64" "mips64" "mips64le" "mips" "mipsle" "ppc64" "ppc64le" "riscv64" "s390x" "wasm"
	do
		target_os_name=${os}
		[ "$os" == "windows" ] && execution_file="curl-go.exe"
		[ "$os" == "darwin" ] && target_os_name="mac"
		
		mkdir -p ../downloads/${version}/${target_os_name}/${arch}
		GOOS=${os} GOARCH=${arch} go build -ldflags "-X VERSION=${version}" -o ../downloads/${version}/${target_os_name}/${arch}/${execution_file} ../curl-go.go 2> /dev/null
		if [ "$?" != "0" ]
		#if compilation failed - remove folders - else copy config file.
		then
		  rm -rf ../downloads/${version}/${target_os_name}/${arch}
		else
		  echo "GOOS=${os} GOARCH=${arch} go build -ldflags "-X VERSION=${version}" -o ../downloads/${version}/${target_os_name}/${arch}/${execution_file} ../curl-go.go"
		fi
	done
done

rsync -avP ../downloads/* root@files.matveynator.ru:/home/files/public_html/curl-go/
