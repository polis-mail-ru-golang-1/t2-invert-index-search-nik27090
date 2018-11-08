package invertIndex

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var sliceFiles2 []File

func init() {
	f := File{Name: "time.txt", Content: ""}
	file, _ := ioutil.ReadFile(f.Name)
	f.Content = string(file)
	sliceFiles2 = append(sliceFiles2, f)
}

func BenchmarkInvertInd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PreInvertIndex(sliceFiles2)
	}
}
