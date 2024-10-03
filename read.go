// Copyright ©2022 The mozcookie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mozcookie

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Read returns the cookies stored in the name file.
func Read(name string) ([]*http.Cookie, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("mozcookie: could not open cookie file %q: %w", name, err)
	}
	defer f.Close()

	return Decode(f)
}

func Decode(r io.Reader) ([]*http.Cookie, error) {
	cs, err := decode(r)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

// // Open returns the cookies stored in the name file, as a cookiejar.
// func Open(name string) (*cookiejar.Jar, error) {
// 	f, err := os.Open(name)
// 	if err != nil {
// 		return nil, fmt.Errorf("mozcookie: could not open cookie file %q: %w", name, err)
// 	}
// 	defer f.Close()
//
// 	cs, err := decode(f)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	cookies := make(map[string][]*http.Cookie, len(cs))
// 	for _, c := range cs {
// 		k := c.Domain
// 		switch {
// 		case c.Secure:
// 			k = "https://" + k
// 		default:
// 			k = "http://" + k
// 		}
// 		cookies[k] = append(cookies[k], c)
// 	}
//
// 	jar, err := cookiejar.New(nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for k, v := range cookies {
// 		u, err := url.Parse(k)
// 		if err != nil {
// 			return nil, fmt.Errorf("mozcookie: could not parse url %q: %w", k, err)
// 		}
// 		jar.SetCookies(u, v)
// 	}
// 	return jar, nil
// }

func decode(r io.Reader) ([]*http.Cookie, error) {
	scan := bufio.NewScanner(r)
	// read header.
	if !scan.Scan() || scan.Err() != nil {
		if scan.Err() != nil {
			return nil, fmt.Errorf("mozcookie: could not read header: %w", scan.Err())
		}
		// empty file, or empty header.
		return nil, nil
	}

	hdr := strings.TrimSpace(scan.Text())
	if !magicRe.MatchString(hdr) {
		return nil, fmt.Errorf("mozcookie: invalid cookie magic")
	}

	var cookies []*http.Cookie
	for scan.Scan() {
		txt := strings.TrimSpace(scan.Text())
		if txt == "" {
			continue
		}

		c := &http.Cookie{
			SameSite: http.SameSiteDefaultMode,
		}

		if strings.HasPrefix(txt, httpOnlyPrefix) {
			c.HttpOnly = true
			txt = txt[len(httpOnlyPrefix):]
		}

		if txt == "" || strings.HasPrefix(txt, "#") {
			continue
		}

		toks := strings.Split(txt, "\t")

		if len(toks) > 0 {
			c.Domain = strings.ToLower(toks[0])
		}
		// domain-specified: toks[1] == "TRUE"
		if len(toks) > 2 {
			c.Path = toks[2]
		}
		if len(toks) > 3 {
			c.Secure = toks[3] == "TRUE"
		}
		if len(toks) > 4 {
			c.RawExpires = toks[4]
		}
		if len(toks) > 5 {
			c.Name = toks[5]
		}
		if len(toks) > 6 {
			c.Value = toks[6]
		}

		if c.RawExpires != "" {
			v, err := strconv.ParseInt(c.RawExpires, 10, 64)
			if err != nil {
				return nil, err
			}
			c.Expires = time.Unix(v, 0).UTC()
		}

		cookies = append(cookies, c)
	}

	if err := scan.Err(); err != nil {
		switch {
		case errors.Is(err, io.EOF):
			// ok
		default:
			return nil, fmt.Errorf("mozcookie: could not decode cookie stream: %w", err)
		}
	}

	return cookies, nil
}
