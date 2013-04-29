/*
 * Copyright (c) 2013 Matt Jibson <matt.jibson@gmail.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package goapp

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"fmt"
	"github.com/mjibson/goon"
	"net/http"
)

func init() {
	http.HandleFunc("/", Main)
}

func Main(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	n := goon.NewGoon(r)

	// key tests

	noid := NoId{}
	if _, err := n.Key(noid); err == nil {
		fmt.Fprintln(w, "expected error on noid")
	}

	var keyTests = []keyTest{
		keyTest{
			HasId{Id: 1},
			datastore.NewKey(c, "HasId", "", 1, nil),
		},
		keyTest{
			HasKind{Id: 1, Kind: "OtherKind"},
			datastore.NewKey(c, "OtherKind", "", 1, nil),
		},
		keyTest{
			HasDefaultKind{Id: 1, Kind: "OtherKind"},
			datastore.NewKey(c, "OtherKind", "", 1, nil),
		},
		keyTest{
			HasDefaultKind{Id: 1},
			datastore.NewKey(c, "DefaultKind", "", 1, nil),
		},
	}

	for _, kt := range keyTests {
		if k, err := n.Key(kt.obj); err != nil {
			fmt.Fprintln(w, "error:", err.Error())
		} else if !k.Equal(kt.key) {
			fmt.Fprintln(w, "keys not equal")
		}
	}

	// datastore tests

	initTest(c)
	if err := n.Get(&HasId{Id: 0}); err == nil {
		fmt.Fprintln(w, "ds: expected error")
	}
	if err := n.Get(&HasId{Id: 1}); err != datastore.ErrNoSuchEntity {
		fmt.Fprintln(w, "ds: expected no such entity")
	}
	// run twice to make sure autocaching works correctly
	if err := n.Get(&HasId{Id: 1}); err != datastore.ErrNoSuchEntity {
		fmt.Fprintln(w, "ds: expected no such entity")
	}

	fmt.Fprintln(w, "done")
}

func initTest(c appengine.Context) {
	keys, _ := datastore.NewQuery("HasId").KeysOnly().GetAll(c, nil)
	datastore.DeleteMulti(c, keys)
	memcache.Flush(c)
}

type keyTest struct {
	obj interface{}
	key *datastore.Key
}

type NoId struct {
}

type HasId struct {
	Id   int64 `datastore:"-" goon:"id"`
	Name string
}

type HasKind struct {
	Id   int64  `datastore:"-" goon:"id"`
	Kind string `datastore:"-" goon:"kind"`
	Name string
}

type HasDefaultKind struct {
	Id   int64  `datastore:"-" goon:"id"`
	Kind string `datastore:"-" goon:"kind,DefaultKind"`
	Name string
}