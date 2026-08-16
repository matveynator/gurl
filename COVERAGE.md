# Code Coverage

The `v1.0.0` release has **90.6% statement coverage** in the main package. The
report was generated with Go 1.26.6 on 2026-08-16:

```text
github.com/matveynator/gurl/gurl.go:27:  String                 100.0%
github.com/matveynator/gurl/gurl.go:31:  Set                    100.0%
github.com/matveynator/gurl/gurl.go:61:  main                     0.0%
github.com/matveynator/gurl/gurl.go:65:  run                    100.0%
github.com/matveynator/gurl/gurl.go:91:  parseOptions           100.0%
github.com/matveynator/gurl/gurl.go:139: execute                 91.3%
github.com/matveynator/gurl/gurl.go:178: buildRequest            81.8%
github.com/matveynator/gurl/gurl.go:228: normalizeURL            90.0%
github.com/matveynator/gurl/gurl.go:245: prepareMultipartForm    86.7%
github.com/matveynator/gurl/gurl.go:269: addMultipartFile        80.0%
github.com/matveynator/gurl/gurl.go:287: writeHeaders            66.7%
github.com/matveynator/gurl/gurl.go:294: writeResponse           93.3%
github.com/matveynator/gurl/gurl.go:322: Error                  100.0%
total:                                  (statements)           90.6%
```

Reproduce the report from the repository root:

```sh
go test -race -count=1 -coverprofile=coverage.out .
go tool cover -func=coverage.out
```

The Security workflow runs the same test suite for every pull request and
rejects coverage below 80%.
