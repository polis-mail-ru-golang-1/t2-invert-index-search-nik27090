package main

import (

)

func main() {
	//срез файлов(название и содержание)


	//инвертированный индекс
	inIn := funcs.InvertIndex(files)

	//срез с файлами в которых поисковая фраза содержиться полностью
	end := funcs.Find(inIn, phrases, files)

	//сортировка файлов по большему кол-ву сопадений
	funcs.SortSearch(end)
}
