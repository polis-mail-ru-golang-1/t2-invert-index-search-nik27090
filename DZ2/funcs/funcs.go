package funcs

import (
	"fmt"
	"strings"
)

type File struct {
	Name    string
	Content string
}

func InvertIndex(sliceFiles []File) map[string]map[string]int {
	InIn := make(map[string]map[string]int)
	for _, f := range sliceFiles {
		sliceStrFile := strings.Split(f.Content, " ")
		for _, word := range sliceStrFile {
			if InIn[word] == nil {
				fileMap := make(map[string]int)
				fileMap[f.Name]++
				InIn[word] = fileMap
			} else {
				InIn[word][f.Name]++
			}
		}
	}
	return InIn
}

func Find(inIn map[string]map[string]int, phrase []string, sliceFiles []File) map[string]int {
	phWords := haveWord(inIn, phrase)
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
func haveWord(inIn map[string]map[string]int, phrase []string) map[string]map[string]int {
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

func SortSearch(endMap map[string]int) {
	if len(endMap) == 0 {
		fmt.Println("Не найденно файлов по данному запросу")
		return
	}
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
		fmt.Println("Файл:", nameFile[i], "; совпадений:", count[i])
	}
}
