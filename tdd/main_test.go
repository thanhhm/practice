package main

import (
	"testing"
)

var mockGetDBID func(id int) int
var mockGetDepsID func() int

type mockDB struct{}

func (mdb mockDB) getDBID(id int) int { return mockGetDBID(id) }

type mockAPIs struct{}

func (mapi mockAPIs) getDepsID() int { return mockGetDepsID() }
func TestGoodFunc2(t *testing.T) {
	testTable := []struct {
		inDB   int
		inAPI  int
		expect bool
	}{
		{1, 1, true}, {1, 2, false}, {0, -1, false},
	}
	mdb := mockDB{}
	mapi := mockAPIs{}

	for _, v := range testTable {
		mockGetDBID = func(id int) int {
			return v.inDB
		}
		mockGetDepsID = func() int {
			return v.inAPI
		}

		if out := goodFunc2(v.inDB, mdb, mapi); out != v.expect {
			t.Fatalf("expect: %t but got %t \n in: %d", v.expect, out, v.inDB)
		}
	}

}
