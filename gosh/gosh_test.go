package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "testing"

func TestGetGoshMapFromArgs(t *testing.T) {
	goshMap, _ := getMap([]string{"foo/bar,git@github.com/MediaMath/foo.git", "salt,git@github.com/MediaMath/salt.git"})

	if len(goshMap) != 2 {
		t.Errorf("|%v|", goshMap)
	}

	if goshMap["foo/bar"] != "git@github.com/MediaMath/foo.git" {
		t.Errorf("|%v|", goshMap)
	}

	if goshMap["salt"] != "git@github.com/MediaMath/salt.git" {
		t.Errorf("|%v|", goshMap)
	}

}

func TestGetGoshMapFromArgsGarbage(t *testing.T) {
	if goshMap, err := getMap([]string{"foo/bar;git@github.com/MediaMath/foo.git", "salt,git@github.com/MediaMath/salt.git"}); err == nil {
		t.Errorf("|%v|", goshMap)
	}
}
