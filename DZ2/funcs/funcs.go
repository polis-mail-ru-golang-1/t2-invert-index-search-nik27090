package funcs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func СreatePhrase() []string {
	fmt.Println("Введите фразу:")
	phrase := ScanStr()
	slicePhrase := strings.Split(phrase, " ")
	return slicePhrase
}

func ScanStr() string {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка ввода:", err)
	}
	return in.Text()
}

var InvertIndex = map[string][]int{}

func CountAndII(sliceStr []string, sliceFiles []FileIndexer) {
	var c int
	for _, vol := range sliceStr {
		index := make([]int, 0)
		j := 0
		for i := range sliceFiles {
			c = strings.Count(sliceFiles[i].Content, vol)
			if c > 0 {
				index = append(index, i)
				j++
				sliceFiles[i].Times++
			}
		}
		if j > 0 {
			InvertIndex[vol] = index
		}
	}
}

type FileIndexer struct {
	Index   int
	Name    string
	Content string
	Times   int
}

func OpenFiles() []FileIndexer {
	fileNames := os.Args[1:]
	sliceFiles := make([]FileIndexer, 0)
	for i, names := range fileNames {
		file, err := ioutil.ReadFile(names)
		Check(err)
		f := FileIndexer{Index: i, Name: names, Content: string(file)}
		sliceFiles = append(sliceFiles, f)
	}
	return sliceFiles
}
