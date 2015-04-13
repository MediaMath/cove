package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import "testing"

func TestGetGoshMapFromArgs(t *testing.T) {
	goshMap, _ := getMap([]string{"foo/bar,git@github.com/MediaMath/foo.git", "salt,git@github.com/MediaMath/salt.git"})

	if len(goshMap) != 2 {
		t.Errorf("Map does not have appropriate number of items |%v|", goshMap)
	}

	if goshMap["foo/bar"] != "git@github.com/MediaMath/foo.git" {
		t.Errorf("Incorrect url for foo/bar", goshMap["foo/bar"])
	}

	if goshMap["salt"] != "git@github.com/MediaMath/salt.git" {
		t.Errorf("Incorrect url for salt", goshMap["salt"])
	}

}

func TestGetGoshMapMixOfImpliedAndExplicit(t *testing.T) {
	goshMap, _ := getMap([]string{"foo/bar/baz", "salt,git@github.com/MediaMath/salt.git"})

	if len(goshMap) != 2 {
		t.Errorf("Map does not have appropriate number of items |%v|", goshMap)
	}

	if goshMap["foo/bar/baz"] != "git@github.com/MediaMath/foo.git" {
		t.Errorf("Incorrect url for foo/bar/baz", goshMap["foo/bar"])
	}

	if goshMap["salt"] != "git@github.com/MediaMath/salt.git" {
		t.Errorf("Incorrect url for salt", goshMap["salt"])
	}
}

func TestImpliedGithubRepo(t *testing.T) {
	implied, _ := impliedGithubRepo("github.com/MediaMath/foo")

	if implied != "git@github.com:MediaMath/foo.git" {
		t.Errorf("Got:%v", implied)
	}

}

func TestImpliedNoHostProducesGithubUrl(t *testing.T) {
	implied, _ := impliedGithubRepo("bar/foo")

	if implied != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", implied)
	}
}

func TestImpliedSubpackage(t *testing.T) {
	impliedWithHost, _ := impliedGithubRepo("github.com/bar/foo/baz")

	if impliedWithHost != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", impliedWithHost)
	}

	impliedNoHost, _ := impliedGithubRepo("bar/foo/baz")

	if impliedNoHost != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", impliedNoHost)
	}

	impliedDeep, _ := impliedGithubRepo("bar/foo/baz/goose/gander")

	if impliedDeep != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", impliedDeep)
	}
}

func TestNonGithubHostFails(t *testing.T) {
	if _, norepo := impliedGithubRepo("foo.com/bar/baz"); norepo == nil {
		t.Errorf("Didnt get error on norepo")
	}
}

func TestSingleParamImplicationFails(t *testing.T) {
	if _, singleParam := impliedGithubRepo("foo"); singleParam == nil {
		t.Errorf("Didnt get error on only host")
	}
}
