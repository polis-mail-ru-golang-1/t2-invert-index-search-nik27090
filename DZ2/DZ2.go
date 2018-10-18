/*
первый вариант
map[string]map[string]int
*/
/*
второй вариант
struct around{
	name string
	count int
}
*/
//сортировка мапы: перевести в слайс по макс значению
package main
//пока что производит поиск не правильно
import (
	"fmt"
	"sort"

	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-nik27090/tree/master/DZ2/funcs"
)

func main() {
	files := funcs.OpenFiles()

	//ввод поисковой фразы и инвертирование в срез
	phrases := funcs.СreatePhrase()

	//интвертированный индекс
	funcs.CountAndII(phrases, files)

	//чем больше слов из фразы встретилось в фале тем он выше, кол-во одного и того же слова не учитывается
	sort.Slice(files, func(i, j int) bool { return files[i].Times > files[j].Times })

	for i := range files {
		if files[i].Times > 0 {
			fmt.Println(files[i].Name, "; совпадений - ", files[i].Times)
		}
	}
}

/*
func scanStr() string {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка ввода:", err)
	}
	return in.Text()
}
*/
/*
type fileIndexer struct {
	index   int
	name    string
	content string
	times   int
}
*/
/*
func openFiles() []funcs.FileIndexer {
	fileNames := os.Args[1:]
	sliceFiles := make([]funcs.FileIndexer, 0)
	for i, names := range fileNames {
		file, err := ioutil.ReadFile(names)
		funcs.Check(err)
		f := funcs.FileIndexer{Index: i, Name: names, Content: string(file)}
		sliceFiles = append(sliceFiles, f)
	}
	return sliceFiles
}
*/
/*func сheck(e error) {
	if e != nil {
		panic(e)
	}
}*/
/*
func createPhrase() []string {
	fmt.Println("Введите фразу:")
	phrase := scanStr()
	slicePhrase := strings.Split(phrase, " ")
	return slicePhrase
}
*/
//var InvertIndex = map[string][]int{}
/*
func countAndII(sliceStr []string, sliceFiles []fileIndexer) {
	var c int
	for _, vol := range sliceStr {
		index := make([]int, 0)
		j := 0
		for i := range sliceFiles {
			c = strings.Count(sliceFiles[i].content, vol)
			if c > 0 {
				index = append(index, i)
				j++
				sliceFiles[i].times++
			}
		}
		if j > 0 {
			invertIndex[vol] = index
		}
	}
}
*/
