package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"./invertIndex"
)

var inIn map[string]map[string]int
var sliceFiles []invertIndex.File

func handler(w http.ResponseWriter, r *http.Request) {

	q := r.FormValue("q")
	if q != "" {
		endMap := find(inIn, q, sliceFiles)
		sortSearch(endMap, w)
	}
}

func main() {
	inIn, sliceFiles = openFiles()

	http.HandleFunc("/search", handler)
	fmt.Println("Server is listening...")
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func openFiles() (map[string]map[string]int, []invertIndex.File) {
	sliceFiles := make([]invertIndex.File, 0)
	directoryFiles := os.Args[1]

	sliceFileInfo, err := ioutil.ReadDir(directoryFiles)
	check(err)

	for i := 0; i < len(sliceFileInfo); i++ {
		dirFile := directoryFiles + "/" + sliceFileInfo[i].Name()
		textFile, err := ioutil.ReadFile(dirFile)
		check(err)
		f := invertIndex.File{Name: sliceFileInfo[i].Name(), Content: string(textFile)}
		sliceFiles = append(sliceFiles, f)
	}
	return invertIndex.PreInvertIndex(sliceFiles), sliceFiles
}

func find(inIn map[string]map[string]int, q string, sliceFiles []invertIndex.File) map[string]int {
	q = strings.ToLower(q)
	phrase := strings.Split(q, " ")
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
func takeGoodFile(inIn map[string]map[string]int, sliceFiles []invertIndex.File, phrase []string) []string {
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

func sortSearch(endMap map[string]int, w http.ResponseWriter) {
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
	for i := 0; i < len(nameFile); i++ {
		fmt.Fprintf(w, "- %s; совпадений - %d\n", nameFile[i], count[i])
	}
}
