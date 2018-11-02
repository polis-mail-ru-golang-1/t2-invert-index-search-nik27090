package invertIndex

import (
	"strings"
	"sync"
)

type File struct {
	Name    string
	Content string
}

func PreInvertIndex(sliceFiles []File) map[string]map[string]int {
	inIn := make(map[string]map[string]int)

	var wg sync.WaitGroup
	wg.Add(len(sliceFiles))

	var mutex sync.Mutex

	for i := 0; i < len(sliceFiles); i++ {
		go invertIndexGo(inIn, sliceFiles[i].Name, sliceFiles[i].Content, &wg, &mutex)
	}

	wg.Wait()
	return inIn
}

func invertIndexGo(inIn map[string]map[string]int, name string, content string, wg *sync.WaitGroup, mutex *sync.Mutex) {
	StrFile := splitTrim(content)
	for _, word := range StrFile {
		mutex.Lock()
		_, ok := inIn[word]
		if !ok {
			fileMap := make(map[string]int)
			fileMap[name]++
			inIn[word] = fileMap
		} else {
			inIn[word][name]++
		}
		mutex.Unlock()
	}
	wg.Done()
}

func splitTrim(in string) []string {
	words := strings.Fields(in)
	for i := 0; i < len(words); i++ {
		words[i] = strings.Trim(words[i], "/.,?!-\"")
		if words[i] == "" {
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
	return words
}
