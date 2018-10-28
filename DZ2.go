package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-nik27090/DZ2/funcs"
)

func main() {
	inIn := make(map[string]map[string]int)

	//срез файлов(название и содержание)
	files := openFiles()

	//ввод фразы с консоли, создание среза
	phrases := createPhrase()

	//инвертированный индекс
	var wg sync.WaitGroup
	wg.Add(len(files))

	ch := make(chan int, 1)
	ch <- 1
	for i := 0; i < len(files); i++ {
		go funcs.InvertIndexGo(inIn, files[i].Name, files[i].Content, &wg, &ch)
	}

	wg.Wait()

	//срез с файлами в которых поисковая фраза содержиться полностью
	end := funcs.Find(inIn, phrases, files)

	//сортировка файлов по большему кол-ву сопадений
	//и вывод результата в консоль
	funcs.SortSearch(end)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func scanStr() string {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка ввода:", err)
	}
	return in.Text()
}

func openFiles() []funcs.File {
	fileNames := os.Args[1:]
	sliceFiles := make([]funcs.File, 0)
	for _, names := range fileNames {
		file, err := ioutil.ReadFile(names)
		check(err)
		f := funcs.File{Name: names, Content: string(file)}
		sliceFiles = append(sliceFiles, f)
	}
	return sliceFiles
}

func createPhrase() []string {
	fmt.Println("Введите фразу:")
	phrase := scanStr()
	slicePhrase := strings.Split(phrase, " ")
	return slicePhrase
}

