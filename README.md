# gurl is a curl-like utility written in pure GO (GOLang), with embedded SSL support and no external libraries.

```
/usr/local/bin/gurl --help
Usage of /usr/local/bin/gurl:
  -head
    	Perform HEAD request.
  -lang string
    	Set Accept-Language header, for example: -lang en-US
  -proxy string
    	Set http/https/socks5 proxy 'type://host:port', example: -proxy 'socks5://127.0.0.1:3128' -proxy 'http://127.0.0.1:8080'
  -timeout duration
    	Set connect and operation timeout. Valid time units are: ns, us or µs, ms, s, m, h. (default 30s)
  -unsafe
    	Disable strict TLS certificate checking.
  -useragent string
    	Set user agent. (default "GURL (https://github.com/matveynator/gurl)")
  -version
    	Output version information.
```

## [↓ Download latest version of gurl.](http://files.matveynator.ru/gurl/latest/) 

- Supported OS: [Linux](http://files.matveynator.ru/gurl/latest/linux), [Windows](http://files.matveynator.ru/gurl/latest/windows), [Android](http://files.matveynator.ru/gurl/latest/android), [Mac](http://files.matveynator.ru/gurl/latest/mac), [IOS](http://files.matveynator.ru/gurl/latest/ios), [FreeBSD](http://files.matveynator.ru/gurl/latest/freebsd), [DragonflyBSD](http://files.matveynator.ru/gurl/latest/dragonfly), [OpenBSD](http://files.matveynator.ru/gurl/latest/openbsd), [NetBSD](http://files.matveynator.ru/gurl/latest/netbsd), [Plan9](http://files.matveynator.ru/gurl/latest/plan9), [AIX](http://files.matveynator.ru/gurl/latest/aix), [Solaris](http://files.matveynator.ru/gurl/latest/solaris), [Illumos](http://files.matveynator.ru/gurl/latest/illumos).
- Supported architectures: x86-32, x86-64, ARM, ARM64, MIPS64, MIPS64le, MIPS, MIPSLE, PPC64, PPC64le, RISCv64, s390x. 


### build gurl yourself 
GOLANG version 1.11 or later is required.
```
git clone https://github.com/matveynator/gurl.git
cd gurl
go build
./gurl https://google.com
```

### how to add new feature?
```
- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Added some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request
```

### get gurl for LINUX/amd64:
```
curl -L 'http://files.matveynator.ru/gurl/latest/linux/amd64/gurl' > /usr/local/bin/gurl; chmod +x /usr/local/bin/gurl;
```

### get gurl for LINUX/386:
```
curl -L 'http://files.matveynator.ru/gurl/latest/linux/386/gurl' > /usr/local/bin/gurl; chmod +x /usr/local/bin/gurl;
```

### get gurl for MAC/amd64:
```
curl -L 'http://files.matveynator.ru/gurl/latest/mac/amd64/gurl' > /usr/local/bin/gurl; chmod +x /usr/local/bin/gurl;
```

### get gurl for MAC/386:
```
curl -L 'http://files.matveynator.ru/gurl/latest/mac/386/gurl' > /usr/local/bin/gurl; chmod +x /usr/local/bin/gurl;
```

## Check all precompiled versions [ here ](http://files.matveynator.ru/gurl/latest/). There are a lot! :) 


[aix/ppc64/gurl](http://files.matveynator.ru/gurl/latest/aix/ppc64/gurl)

[android/arm64/gurl](http://files.matveynator.ru/gurl/latest/android/arm64/gurl)

[dragonfly/amd64/gurl](http://files.matveynator.ru/gurl/latest/dragonfly/amd64/gurl)

[freebsd/386/gurl](http://files.matveynator.ru/gurl/latest/freebsd/386/gurl)

[freebsd/amd64/gurl](http://files.matveynator.ru/gurl/latest/freebsd/amd64/gurl)

[freebsd/arm/gurl](http://files.matveynator.ru/gurl/latest/freebsd/arm/gurl)

[freebsd/arm64/gurl](http://files.matveynator.ru/gurl/latest/freebsd/arm64/gurl)

[illumos/amd64/gurl](http://files.matveynator.ru/gurl/latest/illumos/amd64/gurl)

[ios/amd64/gurl](http://files.matveynator.ru/gurl/latest/ios/amd64/gurl)

[linux/386/gurl](http://files.matveynator.ru/gurl/latest/linux/386/gurl)

[linux/amd64/gurl](http://files.matveynator.ru/gurl/latest/linux/amd64/gurl)

[linux/arm/gurl](http://files.matveynator.ru/gurl/latest/linux/arm/gurl)

[linux/arm64/gurl](http://files.matveynator.ru/gurl/latest/linux/arm64/gurl)

[linux/mips/gurl](http://files.matveynator.ru/gurl/latest/linux/mips/gurl)

[linux/mips64/gurl](http://files.matveynator.ru/gurl/latest/linux/mips64/gurl)

[linux/mips64le/gurl](http://files.matveynator.ru/gurl/latest/linux/mips64le/gurl)

[linux/mipsle/gurl](http://files.matveynator.ru/gurl/latest/linux/mipsle/gurl)

[linux/ppc64/gurl](http://files.matveynator.ru/gurl/latest/linux/ppc64/gurl)

[linux/ppc64le/gurl](http://files.matveynator.ru/gurl/latest/linux/ppc64le/gurl)

[linux/s390x/gurl](http://files.matveynator.ru/gurl/latest/linux/s390x/gurl)

[linux/riscv64/gurl](http://files.matveynator.ru/gurl/latest/linux/riscv64/gurl)

[mac/amd64/gurl](http://files.matveynator.ru/gurl/latest/mac/amd64/gurl)

[mac/arm64/gurl](http://files.matveynator.ru/gurl/latest/mac/arm64/gurl)

[mac/386/gurl](http://files.matveynator.ru/gurl/latest/mac/386/gurl)

[netbsd/386/gurl](http://files.matveynator.ru/gurl/latest/netbsd/386/gurl)

[netbsd/amd64/gurl](http://files.matveynator.ru/gurl/latest/netbsd/amd64/gurl)

[netbsd/arm/gurl](http://files.matveynator.ru/gurl/latest/netbsd/arm/gurl)

[netbsd/arm64/gurl](http://files.matveynator.ru/gurl/latest/netbsd/arm/gurl)

[openbsd/386/gurl](http://files.matveynator.ru/gurl/latest/openbsd/386/gurl)

[openbsd/amd64/gurl](http://files.matveynator.ru/gurl/latest/openbsd/amd64/gurl)

[openbsd/arm/gurl](http://files.matveynator.ru/gurl/latest/openbsd/arm/gurl)

[openbsd/arm64/gurl](http://files.matveynator.ru/gurl/latest/openbsd/arm64/gurl)

[openbsd/mips64/gurl](http://files.matveynator.ru/gurl/latest/openbsd/mips64/gurl)

[plan9/386/gurl](http://files.matveynator.ru/gurl/latest/plan9/386/gurl)

[plan9/amd64/gurl](http://files.matveynator.ru/gurl/latest/plan9/amd64/gurl)

[plan9/arm/gurl](http://files.matveynator.ru/gurl/latest/plan9/arm/gurl)

[solaris/amd64/gurl](http://files.matveynator.ru/gurl/latest/solaris/amd64/gurl)

[windows/386/gurl.exe](http://files.matveynator.ru/gurl/latest/windows/386/gurl.exe)

[windows/amd64/gurl.exe](http://files.matveynator.ru/gurl/latest/windows/amd64/gurl.exe)

[windows/arm/gurl.exe](http://files.matveynator.ru/gurl/latest/windows/arm/gurl.exe)

[windows/arm64/gurl.exe](http://files.matveynator.ru/gurl/latest/windows/arm64/gurl.exe)

and [some other versions...](http://files.matveynator.ru/gurl/latest/)


