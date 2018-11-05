package invertIndex

import (
	"strings"
	"testing"
)

//go test -coverprofile=cover.out

var sliceFiles = make([]File, 0)
var name = "f1.txt"
var content = "Hi, how are you, and you, and you?"
var sliceText []string
var expected map[string]map[string]int

func TestPreInvertIndex(t *testing.T) {

	text := strings.ToLower(content)
	sliceText = strings.Fields(text)

	sliceFiles = append(sliceFiles, File{Name: name, Content: content})
	actual := PreInvertIndex(sliceFiles)

	var m1 = map[string]int{
		name: 1,
	}
	var m2 = map[string]int{
		name: 2,
	}
	var m3 = map[string]int{
		name: 3,
	}
	expected = map[string]map[string]int{
		"hi":  m1,
		"how": m1,
		"are": m1,
		"you": m3,
		"and": m2,
	}

	if !testPreInvertIndex(expected, actual, name) {
		t.Errorf("%v != %v", actual, expected)
	}
}

func testPreInvertIndex(a, b map[string]map[string]int, name string) bool {

	if len(a) != len(b) {
		return false
	}

	for _, vol := range sliceText {
		word := strings.Trim(vol, "().,?!-")
		if a[word][name] != b[word][name] {
			return false
		}
	}

	return true
}
