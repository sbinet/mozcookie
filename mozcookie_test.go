// Copyright ©2022 The ATLAS Collaboration. All rights reserved.
// Use of this source code is governed by the Apache-2
// license that can be found in the LICENSE file.

package mozcookie

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	golden = []*http.Cookie{
		&http.Cookie{Name: "cookie-1", Value: "v$1", Domain: "example.com", Path: "/foo/"},
		&http.Cookie{Name: "cookie-2", Value: "v$2", Domain: ".example.com", Path: "/foo/"},
		&http.Cookie{Name: "cookie-3", Value: "v$3", Domain: "example.com", Path: "/foo/", Expires: time.Unix(1257894000, 0)},
		&http.Cookie{Name: "cookie-4", Value: "v$4", Domain: "example.com", Path: "/foo/", RawExpires: "1257894000"},
		&http.Cookie{Name: "cookie-5", Value: "v$5", Domain: "example.com", Path: "/foo/", RawExpires: "1257894000", Secure: true},
		&http.Cookie{Name: "cookie-6", Value: "v$6", Domain: "example.com", Path: "/foo/", RawExpires: "1257894000", HttpOnly: true},
		&http.Cookie{Name: "cookie-7", Value: "v$7", Domain: ".example.com", Path: "/foo/", RawExpires: "1257894000", HttpOnly: true},
	}
)

// func TestOpen(t *testing.T) {
// 	jar, err := Open("testdata/cookie_golden.txt")
// 	if err != nil {
// 		t.Fatalf("error: %+v", err)
// 	}
//
// 	for _, tc := range []struct {
// 		url  string
// 		want int
// 	}{
// 		{
// 			url:  "http://example.com",
// 			want: 4,
// 		},
// 	} {
// 		t.Run(tc.url, func(t *testing.T) {
// 			u, err := url.Parse(tc.url)
// 			if err != nil {
// 				t.Fatalf("could not parse url: %+v", err)
// 			}
//
// 			cs := jar.Cookies(u)
// 			if got, want := len(cs), tc.want; got != want {
// 				t.Fatalf("invalid number of cookies: got=%d, want=%d", got, want)
// 			}
// 		})
// 	}
// }

func TestRoundTrip(t *testing.T) {
	wbuf := new(bytes.Buffer)
	err := Encode(wbuf, golden)
	if err != nil {
		t.Fatalf("could not encode cookies: %+v", err)
	}

	got, err := Decode(bytes.NewReader(wbuf.Bytes()))
	if err != nil {
		t.Fatalf("could not decode cookies: %+v", err)
	}

	rbuf := new(bytes.Buffer)
	err = Encode(rbuf, got)
	if err != nil {
		t.Fatalf("could not decode cookies: %+v", err)
	}

	if !bytes.Equal(dos2unix(wbuf.Bytes()), dos2unix(rbuf.Bytes())) {
		t.Fatalf("round-trip failed")
	}
}

func TestRead(t *testing.T) {
	cs, err := Read("./testdata/cookie_golden.txt")
	if err != nil {
		t.Fatalf("could not read cookies: %+v", err)
	}

	wbuf := new(bytes.Buffer)
	err = Encode(wbuf, cs)
	if err != nil {
		t.Fatalf("could not encode cookies: %+v", err)
	}

	want, err := os.ReadFile("./testdata/cookie_golden.txt")
	if err != nil {
		t.Fatalf("could not read golden file: %+v", err)
	}

	if !bytes.Equal(dos2unix(wbuf.Bytes()), dos2unix(want)) {
		t.Fatalf("round-trip failed")
	}
}

func TestWrite(t *testing.T) {
	tmp, err := os.MkdirTemp("", "mozcookie-")
	if err != nil {
		t.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "cookie.txt")

	err = Write(fname, golden)
	if err != nil {
		t.Fatalf("could not write cookies: %+v", err)
	}

	got, err := os.ReadFile(fname)
	if err != nil {
		t.Fatalf("could not read output file: %+v", err)
	}

	want, err := os.ReadFile("./testdata/cookie_golden.txt")
	if err != nil {
		t.Fatalf("could not read golden file: %+v", err)
	}

	if !bytes.Equal(dos2unix(got), dos2unix(want)) {
		t.Fatalf("round-trip failed:\ngot:\n%s\n\nwant:\n%s\n", got, want)
	}
}

func dos2unix(b []byte) []byte {
	return bytes.Replace(b, []byte("\r\n"), []byte("\n"), -1)
}
