// Copyright ©2022 The ATLAS Collaboration. All rights reserved.
// Use of this source code is governed by the Apache-2
// license that can be found in the LICENSE file.

package mozcookie_test

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sbinet/mozcookie"
)

func Example() {
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

	fmt.Printf("\ncookies-file:\n=============\n")
	o, err := os.ReadFile("./testdata/cookie.txt")
	if err != nil {
		log.Fatalf("could not read cookies file: %+v", err)
	}
	fmt.Printf("%s=============\n", o)

	// Output:
	// cookies:
	// cookie-1: "v$1"
	// cookie-2: "v$2"
	// cookie-3: "v$3"
	//
	// cookies-file:
	// =============
	// # Netscape HTTP Cookie File
	// # http://curl.haxx.se/rfc/cookie_spec.html
	// # This is a generated file!  Do not edit.
	//
	// golang.org	FALSE	/pkg/	FALSE	0	cookie-1	v$1
	// golang.org	FALSE	/pkg/	TRUE	0	cookie-2	v$2
	// #HttpOnly_golang.org	FALSE	/pkg/	FALSE	0	cookie-3	v$3
	// =============
}
