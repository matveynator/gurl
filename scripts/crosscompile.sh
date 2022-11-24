#!/bin/bash
version="0.1-002"
git_root_path=`git rev-parse --show-toplevel`
execution_file="gurl"

go mod download
go mod vendor

cd ${git_root_path}/scripts;

mkdir -p ${git_root_path}/binaries/${version};

rm -f ${git_root_path}/binaries/latest; 

cd ${git_root_path}/binaries; ln -s ${version} latest; cd ${git_root_path}/scripts;

for os in linux freebsd netbsd openbsd aix android illumos ios solaris plan9 darwin dragonfly windows;
#for os in linux;
do
	for arch in "amd64" "386" "arm" "arm64" "mips64" "mips64le" "mips" "mipsle" "ppc64" "ppc64le" "riscv64" "s390x" "wasm"
	do
		target_os_name=${os}
		[ "$os" == "windows" ] && execution_file="gurl.exe"
		[ "$os" == "darwin" ] && target_os_name="mac"
		
		mkdir -p ../binaries/${version}/${target_os_name}/${arch}

		GOOS=${os} GOARCH=${arch} go build -ldflags "-X main.VERSION=${version}" -o ../binaries/${version}/${target_os_name}/${arch}/${execution_file} ../gurl.go 2> /dev/null
		if [ "$?" != "0" ]
		#if compilation failed - remove folders - else copy config file.
		then
		  rm -rf ../binaries/${version}/${target_os_name}/${arch}
		else
		  echo "GOOS=${os} GOARCH=${arch} go build -ldflags "-X main.VERSION=${version}" -o ../binaries/${version}/${target_os_name}/${arch}/${execution_file} ../gurl.go"
		fi
	done
done

#optional: publish to internet:
rsync -avP ../binaries/* root@files.matveynator.ru:/home/files/public_html/gurl/
