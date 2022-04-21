// Copyright ©2022 The mozcookie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mozcookie

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Write(name string, cookies []*http.Cookie) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("mozcookie: could not create cookie file %q: %w", name, err)
	}
	defer f.Close()

	err = Encode(f, cookies)
	if err != nil {
		return fmt.Errorf("mozcookie: could not encode cookies to file %q: %w", name, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("mozcookie: could not close cookie file %q: %w", name, err)
	}

	return nil
}

func Encode(w io.Writer, cookies []*http.Cookie) error {
	enc := NewEncoder(w)
	for _, c := range cookies {
		err := enc.Encode(c)
		if err != nil {
			return fmt.Errorf("mozcookie: could not write cookie %q: %w", c.Name, err)
		}
	}

	if enc.err != nil {
		return fmt.Errorf("mozcookie: could not encode cookies: %w", enc.err)
	}

	return nil
}

type Encoder struct {
	w io.Writer

	buf  strings.Builder
	err  error
	once sync.Once
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(c *http.Cookie) error {
	enc.once.Do(enc.writeHeader)
	if enc.err != nil {
		return enc.err
	}

	enc.buf.Reset()

	domain := c.Domain
	if c.HttpOnly {
		domain = httpOnlyPrefix + domain
	}
	enc.buf.WriteString(domain)
	enc.buf.WriteByte('\t')

	dot := strings.HasPrefix(c.Domain, ".")
	switch {
	case dot:
		enc.buf.WriteString("TRUE\t")
	default:
		enc.buf.WriteString("FALSE\t")
	}

	enc.buf.WriteString(c.Path)
	enc.buf.WriteByte('\t')

	switch {
	case c.Secure:
		enc.buf.WriteString("TRUE\t")
	default:
		enc.buf.WriteString("FALSE\t")
	}

	switch {
	case c.RawExpires != "":
		enc.buf.WriteString(c.RawExpires)
		enc.buf.WriteByte('\t')
	default:
		switch {
		case c.Expires.IsZero():
			enc.buf.WriteString("0")
		default:
			enc.buf.WriteString(strconv.FormatInt(c.Expires.UTC().Unix(), 10))
		}
		enc.buf.WriteByte('\t')
	}

	enc.buf.WriteString(c.Name)
	enc.buf.WriteByte('\t')

	enc.buf.WriteString(c.Value)
	enc.buf.WriteByte('\n')

	var n int
	n, enc.err = enc.w.Write([]byte(enc.buf.String()))
	if enc.err != nil {
		return enc.err
	}
	if n < enc.buf.Len() {
		enc.err = io.ErrShortWrite
	}
	return enc.err
}

func (enc *Encoder) writeHeader() {
	if enc.err != nil {
		return
	}

	_, enc.err = enc.w.Write([]byte(mozHeader))
}
