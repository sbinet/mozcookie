// Copyright ©2022 The mozcookie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mozcookie is a simple package providing tools to read and write cookies
// from and to files in the Netscape HTTP Cookie File format.
//
// For more informations, see:
//
//	http://curl.haxx.se/rfc/cookie_spec.html
package mozcookie // import "github.com/sbinet/mozcookie"

import (
	"regexp"
)

var (
	magicRe = regexp.MustCompilePOSIX("#( Netscape)? HTTP Cookie File")
)

const (
	httpOnlyAttr   = "HTTPOnly"
	httpOnlyPrefix = "#HttpOnly_"
)

const mozHeader = `# Netscape HTTP Cookie File
# http://curl.haxx.se/rfc/cookie_spec.html
# This is a generated file!  Do not edit.

`
