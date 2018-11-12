package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"./invertIndex"
	"go.uber.org/zap"
)

var sliceFiles []invertIndex.File

type config struct {
	Address string
	Direct  string
}

type result struct {
	Name  string
	Count int
}

type AccessLogger struct {
	ZapLogger *zap.SugaredLogger
}

func searchPage(w http.ResponseWriter, r *http.Request) {

	q := r.FormValue("q")
	zap.S().Info("Search phrase: ", q)
	if q != "" {
		endMap := find(invertIndex.InIn, q, sliceFiles)
		endSlice := sortSearch(endMap, w)
		tmpl := template.Must(template.ParseFiles("html/output.html", "html/dynamicList.html"))
		tmpl.Execute(w, endSlice)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/index.html"))
	tmpl.Execute(w, nil)
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ac.ZapLogger.Info(
			zap.String("URL", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}

func main() {
	//config
	conf, _ := os.Open("config.json")
	defer conf.Close()
	decoder := json.NewDecoder(conf)
	config := config{}
	err := decoder.Decode(&config)
	check(err)

	// zap
	zapLogger, err := zap.NewProduction()
	defer zapLogger.Sync()
	check(err)
	zap.ReplaceGlobals(zapLogger)

	zapLogger.Info("server is started",
		zap.String("address", config.Address),
	)

	AccessLogOut := new(AccessLogger)

	sugar := zapLogger.Sugar().With()
	AccessLogOut.ZapLogger = sugar

	// server stuff
	siteMux := http.NewServeMux()
	siteHandler := AccessLogOut.accessLogMiddleware(siteMux)

	invertIndex.InIn, sliceFiles = openFiles(config.Direct)
	zap.S().Info("InvertIndex built.")

	siteMux.HandleFunc("/search", searchPage)
	siteMux.HandleFunc("/", mainPage)
	http.ListenAndServe(config.Address, siteHandler)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func openFiles(dir string) (map[string]map[string]int, []invertIndex.File) {
	sliceFiles := make([]invertIndex.File, 0)

	sliceFileInfo, err := ioutil.ReadDir(dir)
	check(err)
	for i := 0; i < len(sliceFileInfo); i++ {
		dirFile := dir + "/" + sliceFileInfo[i].Name()
		textFile, err := ioutil.ReadFile(dirFile)
		check(err)
		f := invertIndex.File{Name: sliceFileInfo[i].Name(), Content: string(textFile)}
		sliceFiles = append(sliceFiles, f)
		zap.S().Info("File opened: ", sliceFileInfo[i].Name())
	}
	return invertIndex.PreInvertIndex(sliceFiles), sliceFiles
}

func find(inIn map[string]map[string]int, q string, sliceFiles []invertIndex.File) map[string]int {
	q = strings.ToLower(q)
	phrase := strings.Fields(q)
	for i := 0; i < len(phrase); i++ {
		phrase[i] = strings.Trim(phrase[i], "()/.,?!-\"")
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

func sortSearch(endMap map[string]int, w http.ResponseWriter) []result {
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
	nameCount := make([]result, 0)
	for i := 0; i < len(nameFile); i++ {
		f := result{Name: nameFile[i], Count: count[i]}
		nameCount = append(nameCount, f)
	}
	return nameCount
}

