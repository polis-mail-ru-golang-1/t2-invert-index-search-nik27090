package main

import (
	"testing"

	"./invertIndex"
)

//go test -coverprofile=cover.out

func TestFind(t *testing.T) {
	name := "file.txt"
	var m1 = map[string]int{
		name: 1,
	}
	var m2 = map[string]int{
		name: 2,
	}
	var m3 = map[string]int{
		name: 3,
	}
	inIn = map[string]map[string]int{
		"hi":  m1,
		"how": m1,
		"are": m1,
		"you": m3,
		"and": m2,
	}
	q := "and you"
	var sliceFiles = make([]invertIndex.File, 0)
	sliceFiles = append(sliceFiles, invertIndex.File{Name: name, Content: "Hi, how are you, and you, and you?"})

	actual := find(inIn, q, sliceFiles)

	expected := map[string]int{
		name: 5,
	}

	if !testFind(expected, actual) {
		t.Errorf("%v != %v", actual, expected)
	}
}

func testFind(a map[string]int, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}

	if a["file.txt"] != b["file.txt"] {
		return false
	}
	return true
}
