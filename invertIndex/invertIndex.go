package invertIndex

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var InIn map[string]map[string]int

var SliceFiles []File

type File struct {
	Name    string
	Content string
}

type Result struct {
	Name  string
	Count int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func OpenFiles(dir string) (map[string]map[string]int, []File) {
	sliceFiles := make([]File, 0)

	sliceFileInfo, err := ioutil.ReadDir(dir)
	check(err)

	for i := 0; i < len(sliceFileInfo); i++ {
		dirFile := dir + "/" + sliceFileInfo[i].Name()
		textFile, err := ioutil.ReadFile(dirFile)
		check(err)
		f := File{Name: sliceFileInfo[i].Name(), Content: string(textFile)}
		sliceFiles = append(sliceFiles, f)
		zap.S().Info("File opened: ", sliceFileInfo[i].Name())
	}
	return preInvertIndex(sliceFiles), sliceFiles
}

func Find(inIn map[string]map[string]int, q string, sliceFiles []File) map[string]int {
	q = strings.ToLower(q)
	phrase := strings.Fields(q)
	for i := 0; i < len(phrase); i++ {
		phrase[i] = strings.Trim(phrase[i], "':;()/.,?!-\"")
		if phrase[i] == "" {
			phrase = append(phrase[:i], phrase[i+1:]...)
			i--
		}
	}
	phWords := mapWithQWord(inIn, phrase)
	goodFile := takeGoodFile(inIn, sliceFiles, phrase)
	endMap := make(map[string]int)
	//добавляет число совпадений слов поисковой фразы с текстом файла
	for _, gFile := range goodFile {
		for _, item := range phWords {
			for name, i := range item {
				if gFile == name {
					endMap[name] = endMap[name] + i
				}
			}
		}
	}
	return endMap
}

//уменьшает ИнвИнд до имеющихся слов в поисковой фразе
func mapWithQWord(inIn map[string]map[string]int, phrase []string) map[string]map[string]int {
	phWords := make(map[string]map[string]int)
	for fileWord, _ := range inIn {
		for _, findWord := range phrase {
			if fileWord == findWord {
				phWords[findWord] = inIn[findWord]
			}
		}
	}
	return phWords
}

//создает срез файлов имеющих поисковую фразу полностью
func takeGoodFile(inIn map[string]map[string]int, sliceFiles []File, phrase []string) []string {
	s := 0
	goodFile := make([]string, 0)
	for _, file := range sliceFiles {
		for _, ph := range phrase {
			if _, ok := inIn[ph][file.Name]; ok {
				if inIn[ph][file.Name] != 0 {
					s++
				}
			} else {
				continue
			}
		}
		if s == len(phrase) {
			goodFile = append(goodFile, file.Name)
		}
		s = 0
	}
	return goodFile
}

func SortSearch(endMap map[string]int, w http.ResponseWriter) []Result {
	bufName := ""
	bufCount := 0
	nameFile := make([]string, 0)
	count := make([]int, 0)
	for name, c := range endMap {
		nameFile = append(nameFile, name)
		count = append(count, c)
	}
	for i := 0; i < len(nameFile); i++ {
		for j := i; j < len(nameFile); j++ {
			if count[i] < count[j] {
				bufName = nameFile[i]
				nameFile[i] = nameFile[j]
				nameFile[j] = bufName
				bufCount = count[i]
				count[i] = count[j]
				count[j] = bufCount
			}
		}
	}
	nameCount := make([]Result, 0)
	for i := 0; i < len(nameFile); i++ {
		f := Result{Name: nameFile[i], Count: count[i]}
		nameCount = append(nameCount, f)
	}
	return nameCount
}

func preInvertIndex(sliceFiles []File) map[string]map[string]int {
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
	lowWords := strings.ToLower(in)
	words := strings.Fields(lowWords)
	for i := 0; i < len(words); i++ {
		words[i] = strings.Trim(words[i], "';:()/.,?!-\"")
		if words[i] == "" {
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
	return words
}
