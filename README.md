# mozcookie

[![GoDoc](https://pkg.go.dev/badge/github.com/sbinet/mozcookie?status.svg)](https://pkg.go.dev/github.com/sbinet/mozcookie)

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
