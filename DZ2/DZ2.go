package main

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-nik27090/tree/master/DZ2/funcs"
)

func main() {
	//срез файлов(название и содержание)
	files := funcs.OpenFiles()

	//ввод фразы с консоли, создание среза
	phrases := funcs.СreatePhrase()

	//инвертированный индекс
	inIn := funcs.InvertIndex(files)

	//срез с файлами в которых поисковая фраза содержиться полностью
	end := funcs.Find(inIn, phrases, files)

	//сортировка файлов по большему кол-ву сопадений
	funcs.SortSearch(end)
}
