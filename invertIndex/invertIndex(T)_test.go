package invertIndex

import (
	"strings"
	"testing"
)

var sliceFiles = make([]File, 0)
var name = "f1.txt"
var content = "Hi, how are you, and you, and you?"
var sliceText []string

func TestPreInvertIndex(t *testing.T) {

	strings.ToLower(content)
	sliceText = strings.Split(content, " ")

	sliceFiles = append(sliceFiles, File{Name: name, Content: "Hi, how are you, and you, and you?"})
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
	var expected = map[string]map[string]int{
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
		if a[vol][name] != b[vol][name] {
			return false
		}
	}

	return true
}
