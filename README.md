# mozcookie

[![GitHub release](https://img.shields.io/github/release/sbinet/mozcookie.svg)](https://github.com/sbinet/mozcookie/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/sbinet/mozcookie?status.svg)](https://pkg.go.dev/github.com/sbinet/mozcookie)
[![CI](https://github.com/sbinet/mozcookie/workflows/CI/badge.svg)](https://github.com/sbinet/mozcookie/actions)
[![codecov](https://codecov.io/gh/sbinet/mozcookie/branch/main/graph/badge.svg)](https://codecov.io/gh/sbinet/mozcookie)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/sbinet/mozcookie/blob/main/LICENSE)

`mozcookie` is a simple package providing tools to read and write cookies from and to files in the Netscape HTTP Cookie File format.

For more informations about this format, see:

- [http://curl.haxx.se/rfc/cookie_spec.html](http://curl.haxx.se/rfc/cookie_spec.html)

## Example

```go
	cookies := []*http.Cookie{
		&http.Cookie{Name: "cookie-1", Value: "v$1", Domain: "golang.org", Path: "/pkg/"},
		&http.Cookie{Name: "cookie-2", Value: "v$2", Domain: "golang.org", Path: "/pkg/", Secure: true},
		&http.Cookie{Name: "cookie-3", Value: "v$3", Domain: "golang.org", Path: "/pkg/", HttpOnly: true},
	}
	err := mozcookie.Write("testdata/cookie.txt", cookies)
	if err != nil {
		log.Fatalf("could not write cookies: %+v", err)
	}

	cs, err := mozcookie.Read("testdata/cookie.txt")
	if err != nil {
		log.Fatalf("could not read cookies: %+v", err)
	}

	fmt.Printf("cookies:\n")
	for _, c := range cs {
		fmt.Printf("%s: %q\n", c.Name, c.Value)
	}
```

## License

`mozcookie` is released under the `BSD-3` license.
